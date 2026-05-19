param(
  [string]$OutPath = "docs/screenshots/01-main.png",
  [string]$WindowTitle = "Azeroth Launcher"
)

Add-Type -AssemblyName System.Drawing
Add-Type -AssemblyName System.Windows.Forms
Add-Type @'
using System;
using System.Runtime.InteropServices;
public class Win32 {
  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr hWnd, out RECT lpRect);
  [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr hWnd);
  [DllImport("user32.dll")] public static extern bool ShowWindow(IntPtr hWnd, int nCmdShow);
  [DllImport("user32.dll")] public static extern bool IsIconic(IntPtr hWnd);
  [DllImport("user32.dll")] public static extern bool BringWindowToTop(IntPtr hWnd);
  [DllImport("user32.dll")] public static extern bool SetWindowPos(IntPtr hWnd, IntPtr hWndInsertAfter, int X, int Y, int cx, int cy, uint uFlags);
  [DllImport("dwmapi.dll")] public static extern int DwmGetWindowAttribute(IntPtr hwnd, int dwAttribute, out RECT pvAttribute, int cbAttribute);
  [DllImport("user32.dll")] public static extern bool PrintWindow(IntPtr hWnd, IntPtr hdcBlt, uint nFlags);
  [StructLayout(LayoutKind.Sequential)] public struct RECT { public int Left, Top, Right, Bottom; }
  public const int DWMWA_EXTENDED_FRAME_BOUNDS = 9;
  public const uint PW_RENDERFULLCONTENT = 0x00000002;
  public static readonly IntPtr HWND_TOP = IntPtr.Zero;
  public const uint SWP_NOMOVE = 0x0002;
  public const uint SWP_NOSIZE = 0x0001;
  public const uint SWP_SHOWWINDOW = 0x0040;
}
'@

$proc = Get-Process | Where-Object { $_.MainWindowTitle -eq $WindowTitle -and $_.MainWindowHandle -ne 0 } | Select-Object -First 1
if (-not $proc) { Write-Error "Window '$WindowTitle' not found"; exit 1 }

$h = $proc.MainWindowHandle
if ([Win32]::IsIconic($h)) { [Win32]::ShowWindow($h, 9) | Out-Null }
# Lift to top without stealing focus (SetWindowPos vs SetForegroundWindow).
[Win32]::SetWindowPos($h, [Win32]::HWND_TOP, 0, 0, 0, 0, [Win32]::SWP_NOMOVE -bor [Win32]::SWP_NOSIZE -bor [Win32]::SWP_SHOWWINDOW) | Out-Null
[Win32]::BringWindowToTop($h) | Out-Null
Start-Sleep -Milliseconds 500

$r = New-Object Win32+RECT
$dwmOk = [Win32]::DwmGetWindowAttribute($h, [Win32]::DWMWA_EXTENDED_FRAME_BOUNDS, [ref]$r, [System.Runtime.InteropServices.Marshal]::SizeOf([type][Win32+RECT]))
if ($dwmOk -ne 0) { [Win32]::GetWindowRect($h, [ref]$r) | Out-Null }

$w = $r.Right - $r.Left
$ht = $r.Bottom - $r.Top
if ($w -le 0 -or $ht -le 0) { Write-Error "Invalid window rect: ${w}x${ht}"; exit 1 }

$dir = Split-Path $OutPath
if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }

# Use PrintWindow with PW_RENDERFULLCONTENT — works for layered / DirectComposition
# webview windows that won't render via BitBlt/CopyFromScreen when occluded.
# Note: PrintWindow uses GetWindowRect (full window incl. invisible borders), so
# capture there and crop to DWM bounds afterwards.
$wr = New-Object Win32+RECT
[Win32]::GetWindowRect($h, [ref]$wr) | Out-Null
$ww = $wr.Right - $wr.Left
$wh = $wr.Bottom - $wr.Top

$bmp = New-Object System.Drawing.Bitmap $ww, $wh
$g = [System.Drawing.Graphics]::FromImage($bmp)
$hdc = $g.GetHdc()
$ok = [Win32]::PrintWindow($h, $hdc, [Win32]::PW_RENDERFULLCONTENT)
$g.ReleaseHdc($hdc)
$g.Dispose()

if (-not $ok) { Write-Error "PrintWindow failed"; exit 1 }

# Crop to DWM bounds.
$cropX = $r.Left - $wr.Left
$cropY = $r.Top - $wr.Top
if ($cropX -lt 0) { $cropX = 0 }
if ($cropY -lt 0) { $cropY = 0 }
$cropRect = New-Object System.Drawing.Rectangle($cropX, $cropY, $w, $ht)
$cropped = $bmp.Clone($cropRect, $bmp.PixelFormat)
$bmp.Dispose()

$cropped.Save($OutPath, [System.Drawing.Imaging.ImageFormat]::Png)
$cropped.Dispose()

Write-Host ("Saved {0}  ({1}x{2})" -f $OutPath, $w, $ht)
