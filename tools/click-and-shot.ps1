param(
  [int]$X,                  # window-relative X
  [int]$Y,                  # window-relative Y
  [string]$OutPath,
  [string]$WindowTitle = "Azeroth Launcher",
  [int]$WaitMs = 700
)

Add-Type -AssemblyName System.Drawing
Add-Type -AssemblyName System.Windows.Forms
Add-Type @'
using System;
using System.Runtime.InteropServices;
public class W {
  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr hWnd, out RECT lpRect);
  [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr hWnd);
  [DllImport("user32.dll")] public static extern bool ShowWindow(IntPtr hWnd, int nCmdShow);
  [DllImport("user32.dll")] public static extern bool BringWindowToTop(IntPtr hWnd);
  [DllImport("user32.dll")] public static extern bool SetWindowPos(IntPtr hWnd, IntPtr hWndInsertAfter, int X, int Y, int cx, int cy, uint uFlags);
  [DllImport("user32.dll")] public static extern bool SetCursorPos(int X, int Y);
  [DllImport("user32.dll")] public static extern void mouse_event(uint dwFlags, uint dx, uint dy, uint dwData, UIntPtr dwExtraInfo);
  [DllImport("dwmapi.dll")] public static extern int DwmGetWindowAttribute(IntPtr hwnd, int dwAttribute, out RECT pvAttribute, int cbAttribute);
  [DllImport("user32.dll")] public static extern bool PrintWindow(IntPtr hWnd, IntPtr hdcBlt, uint nFlags);
  [StructLayout(LayoutKind.Sequential)] public struct RECT { public int Left, Top, Right, Bottom; }
  public const uint MOUSEEVENTF_LEFTDOWN = 0x0002;
  public const uint MOUSEEVENTF_LEFTUP = 0x0004;
  public const uint PW_RENDERFULLCONTENT = 0x00000002;
  public const uint SWP_NOMOVE = 0x0002;
  public const uint SWP_NOSIZE = 0x0001;
  public const uint SWP_SHOWWINDOW = 0x0040;
}
'@

$proc = Get-Process | Where-Object { $_.MainWindowTitle -eq $WindowTitle -and $_.MainWindowHandle -ne 0 } | Select-Object -First 1
if (-not $proc) { Write-Error "Window '$WindowTitle' not found"; exit 1 }

$h = $proc.MainWindowHandle
[W]::ShowWindow($h, 9) | Out-Null
[W]::SetWindowPos($h, [IntPtr]::Zero, 0, 0, 0, 0, [W]::SWP_NOMOVE -bor [W]::SWP_NOSIZE -bor [W]::SWP_SHOWWINDOW) | Out-Null
[W]::BringWindowToTop($h) | Out-Null
[W]::SetForegroundWindow($h) | Out-Null
Start-Sleep -Milliseconds 500

$r = New-Object W+RECT
$dwmOk = [W]::DwmGetWindowAttribute($h, 9, [ref]$r, [System.Runtime.InteropServices.Marshal]::SizeOf([type][W+RECT]))
if ($dwmOk -ne 0) { [W]::GetWindowRect($h, [ref]$r) | Out-Null }

$absX = $r.Left + $X
$absY = $r.Top + $Y
[W]::SetCursorPos($absX, $absY) | Out-Null
Start-Sleep -Milliseconds 100
[W]::mouse_event([W]::MOUSEEVENTF_LEFTDOWN, 0, 0, 0, [UIntPtr]::Zero)
Start-Sleep -Milliseconds 50
[W]::mouse_event([W]::MOUSEEVENTF_LEFTUP, 0, 0, 0, [UIntPtr]::Zero)

Start-Sleep -Milliseconds $WaitMs

# Re-fetch handle in case modal opened a new window.
$proc2 = Get-Process | Where-Object { $_.MainWindowTitle -eq $WindowTitle -and $_.MainWindowHandle -ne 0 } | Select-Object -First 1
$h2 = if ($proc2) { $proc2.MainWindowHandle } else { $h }
[W]::SetWindowPos($h2, [IntPtr]::Zero, 0, 0, 0, 0, [W]::SWP_NOMOVE -bor [W]::SWP_NOSIZE -bor [W]::SWP_SHOWWINDOW) | Out-Null
Start-Sleep -Milliseconds 200

$r2 = New-Object W+RECT
$dwmOk = [W]::DwmGetWindowAttribute($h2, 9, [ref]$r2, [System.Runtime.InteropServices.Marshal]::SizeOf([type][W+RECT]))
if ($dwmOk -ne 0) { [W]::GetWindowRect($h2, [ref]$r2) | Out-Null }
$w = $r2.Right - $r2.Left
$ht = $r2.Bottom - $r2.Top

$wr = New-Object W+RECT
[W]::GetWindowRect($h2, [ref]$wr) | Out-Null
$ww = $wr.Right - $wr.Left
$wh = $wr.Bottom - $wr.Top

$bmp = New-Object System.Drawing.Bitmap $ww, $wh
$g = [System.Drawing.Graphics]::FromImage($bmp)
$hdc = $g.GetHdc()
$ok = [W]::PrintWindow($h2, $hdc, [W]::PW_RENDERFULLCONTENT)
$g.ReleaseHdc($hdc)
$g.Dispose()
if (-not $ok) { Write-Error "PrintWindow failed"; exit 1 }

$cropX = $r2.Left - $wr.Left
$cropY = $r2.Top - $wr.Top
if ($cropX -lt 0) { $cropX = 0 }
if ($cropY -lt 0) { $cropY = 0 }
$cropRect = New-Object System.Drawing.Rectangle($cropX, $cropY, $w, $ht)
$cropped = $bmp.Clone($cropRect, $bmp.PixelFormat)
$bmp.Dispose()

$dir = Split-Path $OutPath
if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }

$cropped.Save($OutPath, [System.Drawing.Imaging.ImageFormat]::Png)
$cropped.Dispose()
Write-Host ("Clicked ({0},{1}) -> {2}  ({3}x{4})" -f $X, $Y, $OutPath, $w, $ht)
