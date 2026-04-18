```markdown
# Design System Specification: The Cognitive Command

## 1. Overview & Creative North Star
**Creative North Star: "The Cognitive Laboratory"**
This design system moves away from the cluttered, line-heavy aesthetics of traditional IDEs. Instead, it adopts the persona of a high-end, editorial laboratory where human intent and machine intelligence converge. We treat the interface not as a grid of containers, but as a fluid "Cognitive Workspace." 

The system achieves a premium feel through **Intentional Asymmetry** and **Optical Depth**. By eschewing standard 1px borders in favor of tonal shifts and liquid glass overlays, we create an environment that feels expansive yet surgically precise. Every element exists to facilitate "The Flow," where the agent's thoughts and the user's commands are indistinguishable.

---

## 2. Colors & Surface Philosophy
The palette is rooted in deep obsidian tones, punctuated by a luminous "Primary Pulse" that signifies AI activity.

### The "No-Line" Rule
Standard UI relies on borders to define space. **This system prohibits 1px solid borders for sectioning.** 
- Boundaries must be defined solely through background color shifts.
- Use `surface-container-low` (#1B1B20) for the main work area and `surface-container-highest` (#35343a) for active contextual panels.
- This creates a "seamless" transition that feels modern and bespoke.

### Surface Hierarchy & Nesting
Treat the UI as a physical stack of frosted obsidian sheets. 
- **Base Layer**: `surface` (#131318).
- **Secondary Workspace**: `surface-container-low` (#1B1B20).
- **Interactive Modals/Overlays**: Use "Liquid Glass"—a combination of `surface-variant` (#35343a) at 60% opacity with a `20px` backdrop-blur.

### The "Glass & Gradient" Rule
To inject "soul" into a technical environment, CTAs and Hero states should utilize subtle chromatic gradients.
- **Primary CTA**: A linear gradient from `primary_container` (#5B5FFF) to `inverse_primary` (#4345e8) at a 135-degree angle.
- **Success State**: A soft radial glow using `success` (#4ADE80) at 10% opacity behind the component to signify "System Health."

---

## 3. Typography
We utilize a high-contrast typographic scale to separate "Interface Instruction" from "Machine Output."

- **The UI Voice (Inter)**: Used for all navigational and structural elements. 
    - *Editorial Note*: Use `display-sm` (2.25rem) with negative letter-spacing (-0.02em) for high-level status headers to create an authoritative, magazine-like feel.
- **The Technical Truth (JetBrains Mono)**: Reserved strictly for code, terminal outputs, and raw data strings. 
    - *Styling*: Use `body-sm` (0.75rem) for code to maximize information density while maintaining legibility.

| Level | Token | Font | Size | Weight |
| :--- | :--- | :--- | :--- | :--- |
| Display | `display-md` | Inter | 2.75rem | 700 |
| Headline | `headline-sm` | Inter | 1.5rem | 600 |
| Code | `body-md` | JetBrains Mono | 0.875rem | 400 |
| Label | `label-sm` | Inter | 0.6875rem | 500 (Uppercase) |

---

## 4. Elevation & Depth
Depth is achieved through **Tonal Layering** rather than structural scaffolding.

- **The Layering Principle**: Instead of shadows, nest containers. A `surface-container-highest` card sitting on a `surface-container-low` background provides a natural, sophisticated lift.
- **Ambient Shadows**: For floating "Agent" elements, use a highly diffused shadow: `0px 24px 48px rgba(0, 0, 0, 0.4)`. The shadow should feel like a soft cloud, not a hard drop.
- **The "Ghost Border" Fallback**: If separation is technically required for accessibility, use a "Ghost Border": `outline-variant` (#454555) at 15% opacity. This provides a "suggestion" of a boundary without breaking the liquid aesthetic.

---

## 5. Components

### 5.1 Buttons & Inputs
- **Primary Button**: 6px border radius (`lg`). Background: `primary_container`. Text: `on_primary_container`. No border. On hover, increase `surface-tint` luminosity.
- **Input Fields**: Use `surface-container-lowest` (#0E0E13). The focus state should not be a border change, but a 2px bottom-accent in `primary` (#C0C1FF).

### 5.2 Conversational "Thought Bubbles"
These represent the AI’s internal processing.
- **Visuals**: Asymmetric 6px radius. Background: `surface-container-high` (#2A292F). 
- **Detail**: A subtle 1px "Ghost Border" on the top and left edges only to catch the "light."

### 5.3 Status Cards & Pulse Animations
- **The Pulse**: For active AI tasks, use a `primary` (#C0C1FF) 2px dot with an outer concentric ring that scales and fades (0% to 40% opacity) every 2 seconds.
- **Layout**: No internal dividers. Use `label-sm` for metadata headers and `title-md` for the primary status.

### 5.4 Dual-Pane Diff Views
- **Logic**: Forbid the use of harsh red/green blocks. 
- **Styling**: Use `error_container` (#93000A) at 20% opacity for deletions and `primary_container` (#5B5FFF) at 15% opacity for additions. This keeps the IDE's professional dark-mode integrity intact.

### 5.5 Gantt-Style Timelines
- **Visuals**: Timeline tracks use `surface-container-highest`. 
- **Active Tasks**: Rounded pill shapes using `primary_fixed_dim`. 
- **Connection Lines**: Use `outline` (#908FA1) at 30% opacity with a `1px` dash-array for a "blueprint" feel.

---

## 6. Do's and Don'ts

### Do:
- **Do** use `vertical-align: middle` for all icon-label pairings to maintain the "Surgical" look.
- **Do** utilize `8px` increments for all padding and margins to ensure the 8px grid system feels intentional.
- **Do** use `surface-bright` for hover states on dark backgrounds to create a "lifting" effect.

### Don't:
- **Don't** use pure black (#000000) or pure white (#FFFFFF). Always use the themed neutrals (`surface` and `on_surface`).
- **Don't** use 1px solid, 100% opaque borders. This is the fastest way to make the design look like a generic template.
- **Don't** use standard "Blue" for success. Only use the `success` (#4ADE80) token to ensure the AI's "Health" status is clear.
- **Don't** crowd the interface. If a screen feels full, increase the `surface-container` nesting depth rather than adding more lines.

---

## 7. Motion & Interaction
- **Transition**: All state changes (hover, focus) must use a `200ms cubic-bezier(0.4, 0, 0.2, 1)` easing.
- **The Liquid Entrance**: Overlays and "Liquid Glass" components should slide in with a slight `5px` Y-axis offset and an opacity fade, mimicking the movement of a physical glass pane being placed on a desk.```