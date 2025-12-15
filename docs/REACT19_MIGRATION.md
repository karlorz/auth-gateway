# React 19 + Semi UI Migration Plan

> Migration guide for upgrading Auth Gateway frontend from React 18 to React 19 with Semi UI v2.88.2

## Overview

### Semi Design vs Semi UI

**They are the same project:**
- `semi-design` = the GitHub repo name (`DouyinFE/semi-design`)
- `@douyinfe/semi-ui` = the npm package name for React 18
- `@douyinfe/semi-ui-19` = the npm package name for React 19

### Current State

| Package | Current Version | Target Version |
|---------|-----------------|----------------|
| react | ^18.2.0 | ^19.0.0 |
| react-dom | ^18.2.0 | ^19.0.0 |
| @douyinfe/semi-ui | ^2.69.1 | Remove |
| @douyinfe/semi-icons | ^2.63.1 | Keep or update |
| @douyinfe/semi-ui-19 | N/A | ^2.88.2 |

## Why Two Packages?

Semi maintains **two separate packages** because React 19 removed critical APIs:

| Removed API | Impact |
|-------------|--------|
| `ReactDOM.render()` | Used internally by Toast, Notification, Modal |
| `ReactDOM.findDOMNode()` | Used internally by Tooltip, Popover, Dropdown |

The `@douyinfe/semi-ui-19` package has completely reimplemented these components to work without the removed APIs.

| React Version | Package to Use |
|---------------|----------------|
| React 18.x | `@douyinfe/semi-ui` |
| React 19.x | `@douyinfe/semi-ui-19` |

Both packages will receive updates for at least one year.

## React 19 Breaking Changes

### Removed APIs

| Removed API | Replacement |
|-------------|-------------|
| `ReactDOM.render()` | `createRoot().render()` |
| `ReactDOM.hydrate()` | `hydrateRoot()` |
| `ReactDOM.findDOMNode()` | Use `useRef` / DOM refs |
| String refs (`ref="input"`) | Callback refs or `useRef` |
| `PropTypes` | TypeScript |

### Behavioral Changes

1. **Error handling**: Errors no longer re-throw automatically
   - Uncaught errors report to `window.reportError`
   - Caught errors go to `console.error`

2. **StrictMode improvements**: Memoization functions (`useMemo`, `useCallback`) now reuse results between double renders in development mode

3. **Suspense optimization**: Fallbacks commit immediately rather than waiting for sibling trees

### TypeScript Changes

- `useRef` now mandates an argument
- `ReactElement` props default to `unknown` instead of `any`
- Global `JSX` namespace moved to `React.JSX`

## Migration Steps

### Step 1: Update React

```bash
cd web
bun add react@19 react-dom@19
```

### Step 2: Replace Semi UI Package

```bash
bun remove @douyinfe/semi-ui
bun add @douyinfe/semi-ui-19@2.88.2
```

### Step 3: Update All Imports

Replace all imports throughout the codebase:

```javascript
// Before
import { Button, Toast, Modal } from '@douyinfe/semi-ui';

// After
import { Button, Toast, Modal } from '@douyinfe/semi-ui-19';
```

**Files to update:**
- All files in `web/src/` that import from `@douyinfe/semi-ui`

You can use this command to find all files:
```bash
grep -r "@douyinfe/semi-ui" web/src --include="*.jsx" --include="*.js" -l
```

### Step 4: Fix Class Component Wrappers

Components like `Tooltip`, `Popover`, `PopConfirm`, `Dropdown` may not work correctly if their children are class components.

**Workaround:** Wrap class components in a real DOM node:

```jsx
// Before (may break in React 19)
<Tooltip content="tip">
  <MyClassComponent />
</Tooltip>

// After (works)
<Tooltip content="tip">
  <span><MyClassComponent /></span>
</Tooltip>
```

### Step 5: Update Entry Point (if needed)

Check `web/src/index.jsx` for legacy ReactDOM usage:

```javascript
// Before (React 18)
import ReactDOM from 'react-dom';
ReactDOM.render(<App />, document.getElementById('root'));

// After (React 19)
import { createRoot } from 'react-dom/client';
const root = createRoot(document.getElementById('root'));
root.render(<App />);
```

### Step 6: Test Affected Components

Components with changed implementations that need extra testing:

| Component | Risk Level | Notes |
|-----------|------------|-------|
| Modal | High | Reimplemented without ReactDOM.render |
| Toast | High | Reimplemented without ReactDOM.render |
| Notification | High | Reimplemented without ReactDOM.render |
| Tooltip | Medium | May break with class component children |
| Popover | Medium | May break with class component children |
| PopConfirm | Medium | May break with class component children |
| Dropdown | Medium | May break with class component children |

## Verification Checklist

- [ ] React and react-dom updated to v19
- [ ] @douyinfe/semi-ui removed
- [ ] @douyinfe/semi-ui-19 installed
- [ ] All imports updated from `semi-ui` to `semi-ui-19`
- [ ] Entry point uses `createRoot` instead of `ReactDOM.render`
- [ ] No `findDOMNode` deprecation warnings in console
- [ ] No `active` boolean attribute warnings (should be fixed in v2.88.2)
- [ ] Toast notifications work correctly
- [ ] Modal dialogs work correctly
- [ ] Tooltips display properly
- [ ] Dropdown menus function correctly
- [ ] All auth flows tested (login, register, 2FA, passkey)

## Rollback Plan

If issues arise, rollback to React 18:

```bash
cd web
bun remove @douyinfe/semi-ui-19
bun add react@18.2.0 react-dom@18.2.0 @douyinfe/semi-ui@2.88.2
# Then revert import changes
```

## References

- [Semi Design React 19 Guide](https://semi.design/en-US/start/react19)
- [React 19 Upgrade Guide](https://react.dev/blog/2024/04/25/react-19-upgrade-guide)
- [React v19 Announcement](https://react.dev/blog/2024/12/05/react-19)
- [Semi Design GitHub Releases](https://github.com/DouyinFE/semi-design/releases)
- [Semi Design Issue #2615 - findDOMNode](https://github.com/DouyinFE/semi-design/issues/2615)

## Version History

| Date | Version | Changes |
|------|---------|---------|
| 2025-12-15 | 1.0 | Initial migration plan |
