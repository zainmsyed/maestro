# Components

<!-- source: seed -->
<!-- Living registry of established components. Populated incrementally by UI stories. -->

## App Shell
<!-- source: story-005 -->
- Global SPA layout uses a narrow left icon sidebar, a 52px topbar, a flexible main content canvas, and a 28px status bar.
- Sidebar: dark `--bg2`, right border, Fraunces “M” logo, square icon nav buttons with soft radius. Active state uses `--bg4` and gold text.
- Topbar: page title in Fraunces, project/status metadata in DM Mono, right-aligned segmented view tabs and primary import/sync action.
- Main canvas: roomy padding, no actual screen components yet; story-005 may show neutral placeholder panels only.
- Status bar: compact operational metadata row with muted labels and colored status dots.

## Buttons and Tabs
<!-- source: story-005 -->
- Buttons are mono, small uppercase-ish tracking, transparent by default, thin border, 5–6px radius.
- Primary action uses gold text and darker gold border; hover fills with subtle gold wash.
- Segmented tabs sit on `--bg3` with inner active state on `--bg` and a thin border.

## Cards / Panels
<!-- source: story-005 -->
- Cards use layered dark backgrounds, 1px translucent border, 8px radius, and 18–24px padding.
- Metric-like values may use Fraunces light weight; labels use small uppercase mono text.
