# New22CodexHome Specification

## Overview

- Target: remote static `/volume2/docker/new-api/custom-home/index.html`
- Reference: `https://new.22codex.xyz/`
- Reference screenshot: `C:\Users\yangyang\AppData\Local\Temp\codex-clipboard-b110b7f5-e032-40be-ac92-7b365fd9fef6.png`
- Interaction model: fixed host navigation + static custom-home content + CSS decorative motion

## Foundation

- Viewport baseline: `1440 × 1000`
- Host navigation: `64px` fixed header
- Reference page body: `#1b1918`-class near-black warm background
- Body font: `"Public Sans", sans-serif`
- Foreground: warm white
- Palette: amber `#ffc96f`, pink `#ff9dbb`, cyan `#84f7e3`, muted sage/blue accents
- Custom homepage source: `/api/home_page_content`, fully self-contained inline CSS/HTML
- External media dependencies: none
- Inline SVG count: `1`
- Section count: `8`

## Hero

- Desktop layout: two columns inside a centered `max-width: 1180px`-class container
- Top offset: approximately `116px`, clearing the fixed host navigation
- Left card: warm translucent brown surface, thin pale border, `~30px` radius
- Right card: black-brown control panel with cyan/pink/amber atmospheric gradients
- Heading: two lines, heavy display weight; second line uses amber → pink → cyan gradient
- Endpoint panel: single service URL, monospace bold text
- Capability strip: three dark mini-panels
- Supporting list: three cyan gradient square bullets
- Actions: gradient primary button and outlined secondary button
- Bottom tags: four compact status chips
- Right visual: top chips, CSS prisms, elliptical tracks, central glass card, status labels, three metric cards

## Stats rail

- Four equal columns at desktop
- Warm translucent background with 1px separators
- Large values: `50+`, `100+`, `Smart Route`, `Live Trace`
- Accent underline: amber/pink/cyan gradient

## Remaining sections

- Feature engine: asymmetric six-card grid
- Launch flow: three numbered cards
- Tool matrix: inline chips; retain Y² NewAPI external tool links
- Control room: dark dashboard panel
- FAQ: two-column accordion-style cards
- CTA: dark rounded banner with primary and secondary links

## Customization contract

- `小安` branding becomes `Y² NewAPI`
- API base becomes `https://yynewapi.yangyangnj.top/v1`
- Internal links target the top-level NewAPI window
- Preserve current Y² tool links
- Keep “自助购卡 1 号”
- Do not restore “自助购卡 2 号”

## Responsive behavior

- `1440px`: exact two-column hero and four-column stats rail
- `768px`: narrower two-column or stacked transition according to reference CSS
- `390px`: single-column hero, wrapped stats/tools, no horizontal overflow
