export function isSmallViewport(): boolean {
  // 480 is breakpoint for tabs disappearing
  return window.innerWidth < 480
}