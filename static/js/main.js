"use strict";
(() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropSymbols = Object.getOwnPropertySymbols;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __propIsEnum = Object.prototype.propertyIsEnumerable;
  var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
  var __spreadValues = (a, b) => {
    for (var prop in b || (b = {}))
      if (__hasOwnProp.call(b, prop))
        __defNormalProp(a, prop, b[prop]);
    if (__getOwnPropSymbols)
      for (var prop of __getOwnPropSymbols(b)) {
        if (__propIsEnum.call(b, prop))
          __defNormalProp(a, prop, b[prop]);
      }
    return a;
  };

  // typescript/utils/platform.ts
  var isMatchMediaChangeEventListenerSupported = typeof MediaQueryList === "function" && typeof MediaQueryList.prototype.addEventListener === "function";

  // typescript/color-scheme.ts
  var setColorSchemeAttribute = (value) => {
    document.documentElement.setAttribute("data-color-scheme", value);
  };
  var setColorSchemeMeta = (value) => {
    const meta = document.head.querySelector('meta[name="color-scheme"]');
    if (meta) {
      meta.content = value;
    }
  };
  var lightScheme = matchMedia("(prefers-color-scheme: light)");
  var handleColorScheme = () => {
    if (lightScheme.matches) {
      setColorSchemeAttribute("light");
    } else {
      setColorSchemeAttribute("dark");
    }
  };
  function InitalizeColorScheme(colorScheme) {
    switch (colorScheme) {
      case "follow-system": {
        handleColorScheme();
        setColorSchemeMeta("dark light");
        if (isMatchMediaChangeEventListenerSupported) {
          lightScheme.addEventListener("change", handleColorScheme);
        } else {
          lightScheme.addListener(handleColorScheme);
        }
        break;
      }
      default:
        setColorSchemeAttribute(colorScheme);
        setColorSchemeMeta(colorScheme);
        if (isMatchMediaChangeEventListenerSupported) {
          lightScheme.removeEventListener("change", handleColorScheme);
        } else {
          lightScheme.removeListener(handleColorScheme);
        }
        break;
    }
  }

  // typescript/share-button.ts
  function ShareButton() {
    const i = document.getElementById("share");
    const shareButton = document.getElementById("btn-share");
    shareButton && shareButton.addEventListener("click", () => {
      i.select();
      document.execCommand("copy");
      i.blur();
      shareButton.classList.add("copied");
    });
  }

  // typescript/utils/dom.ts
  var readyStateListeners = new Set();
  var isDOMReady = () => document.readyState === "complete" || document.readyState === "interactive";
  var addDOMReadyListener = (listener) => readyStateListeners.add(listener);
  if (!isDOMReady()) {
    const onReadyStateChange = () => {
      if (isDOMReady()) {
        document.removeEventListener("readystatechange", onReadyStateChange);
        readyStateListeners.forEach((listener) => listener());
        readyStateListeners.clear();
      }
    };
    document.addEventListener("readystatechange", onReadyStateChange);
  }
  var removeElement = (element) => {
    element && element.remove();
  };

  // typescript/third-party.ts
  var fillInformationOnForm = (key, value) => {
    if (!value) {
      return;
    }
    const element = document.body.querySelector(`form [name="${key}"`);
    if (!element) {
      return;
    }
    element.value = value;
  };
  var DEFAULT_USERSTYLE_META = `/* ==UserStyle==
@name           A new userstyle!
@namespace      userstyles.world
@version        1.0.0
==/UserStyle== */
`;
  var handleSourceCode = (sourceCode) => {
    if (!sourceCode) {
      return;
    }
    if (!/\/\* *?==UserStyle==/g.test(sourceCode)) {
      sourceCode = `${DEFAULT_USERSTYLE_META}${sourceCode}`;
    }
    fillInformationOnForm("code", sourceCode);
  };
  var onMessage = (ev) => {
    if (!ev.data || !ev.data.type) {
      return;
    }
    const { type, data } = ev.data;
    switch (type) {
      case "usw-remove-stylus-button": {
        removeElement(document.querySelector("a#stylus"));
      }
      case "usw-fill-new-style": {
        if (window.location.pathname !== "/api/oauth/authorize_style/new" || !data) {
          return;
        }
        fillInformationOnForm("name", data["name"]);
        const metaData = data["metadata"];
        if (metaData) {
          fillInformationOnForm("description", metaData["description"]);
          fillInformationOnForm("license", metaData["license"]);
          fillInformationOnForm("homepage", metaData["license"]);
        }
        handleSourceCode(data["sourceCode"]);
      }
    }
  };
  window.addEventListener("message", onMessage);
  function BroadcastReady() {
    window.dispatchEvent(new MessageEvent("message", {
      data: { type: "usw-ready" },
      origin: "https://userstyles.world"
    }));
  }

  // typescript/utils/storage.ts
  var DEFAULT_SETTINGS = {
    colorScheme: "follow-system"
  };
  var localStorageKey = "user-preferences";
  var settings = null;
  function getSettings() {
    if (settings) {
      return settings;
    }
    const MaybeSettings = localStorage.getItem(localStorageKey);
    if (!MaybeSettings) {
      localStorage.setItem(localStorageKey, JSON.stringify(DEFAULT_SETTINGS));
      return DEFAULT_SETTINGS;
    }
    const savedSettings = getValidatedObject(JSON.parse(MaybeSettings), DEFAULT_SETTINGS);
    settings = __spreadValues(__spreadValues({}, DEFAULT_SETTINGS), savedSettings);
    localStorage.setItem(localStorageKey, JSON.stringify(settings));
    return settings;
  }
  function storeNewSettings(newSettings) {
    settings = __spreadValues(__spreadValues({}, settings), newSettings);
    localStorage.setItem(localStorageKey, JSON.stringify(settings));
  }
  function getValidatedObject(source, compare) {
    const result = {};
    if (source == null || typeof source !== "object" || Array.isArray(source)) {
      return null;
    }
    Object.keys(source).forEach((key) => {
      const value = source[key];
      if (value == null || compare[key] == null) {
        return;
      }
      const array1 = Array.isArray(value);
      const array2 = Array.isArray(compare[key]);
      if (array1 || array2) {
        if (array1 && array2) {
          result[key] = value;
        }
      } else if (typeof value === "object" && typeof compare[key] === "object") {
        result[key] = getValidatedObject(value, compare[key]);
      } else if (typeof value === typeof compare[key]) {
        result[key] = value;
      }
    });
    return result;
  }

  // typescript/user-settings.ts
  var PREFIX = "usr-settings";
  function SetValues(settings2) {
    if (!window.location.pathname.startsWith("/account")) {
      return;
    }
    document.getElementById(`${PREFIX}--color-scheme`).value = settings2.colorScheme;
  }
  function SaveUserSettingsButton(onSettingsUpdate2) {
    const saveButton = document.getElementById(`${PREFIX}--save`);
    saveButton && saveButton.addEventListener("click", () => {
      const newSettings = {};
      newSettings.colorScheme = document.getElementById(`${PREFIX}--color-scheme`).value;
      storeNewSettings(newSettings);
      onSettingsUpdate2();
    });
  }

  // typescript/main.ts
  var WhenDOMReady = () => {
    ShareButton();
    BroadcastReady();
    SaveUserSettingsButton(onSettingsUpdate);
    SetValues(getSettings());
  };
  if (isDOMReady()) {
    WhenDOMReady();
  } else {
    addDOMReadyListener(WhenDOMReady);
  }
  var onSettingsUpdate = () => {
    const settings2 = getSettings();
    InitalizeColorScheme(settings2.colorScheme);
  };
  InitalizeColorScheme(getSettings().colorScheme);
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vdHlwZXNjcmlwdC91dGlscy9wbGF0Zm9ybS50cyIsICIuLi8uLi90eXBlc2NyaXB0L2NvbG9yLXNjaGVtZS50cyIsICIuLi8uLi90eXBlc2NyaXB0L3NoYXJlLWJ1dHRvbi50cyIsICIuLi8uLi90eXBlc2NyaXB0L3V0aWxzL2RvbS50cyIsICIuLi8uLi90eXBlc2NyaXB0L3RoaXJkLXBhcnR5LnRzIiwgIi4uLy4uL3R5cGVzY3JpcHQvdXRpbHMvc3RvcmFnZS50cyIsICIuLi8uLi90eXBlc2NyaXB0L3VzZXItc2V0dGluZ3MudHMiLCAiLi4vLi4vdHlwZXNjcmlwdC9tYWluLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJleHBvcnQgY29uc3QgaXNNYXRjaE1lZGlhQ2hhbmdlRXZlbnRMaXN0ZW5lclN1cHBvcnRlZCA9IChcbiAgICAnZnVuY3Rpb24nID09PSB0eXBlb2YgTWVkaWFRdWVyeUxpc3QgJiZcbiAgICAnZnVuY3Rpb24nID09PSB0eXBlb2YgTWVkaWFRdWVyeUxpc3QucHJvdG90eXBlLmFkZEV2ZW50TGlzdGVuZXJcbik7XG4iLCAiaW1wb3J0IHtpc01hdGNoTWVkaWFDaGFuZ2VFdmVudExpc3RlbmVyU3VwcG9ydGVkfSBmcm9tICcuL3V0aWxzL3BsYXRmb3JtJztcbmltcG9ydCB0eXBlIHtVc2VyU2V0dGluZ3N9IGZyb20gJy4vdXRpbHMvc3RvcmFnZSc7XG5cbmNvbnN0IHNldENvbG9yU2NoZW1lQXR0cmlidXRlID0gKHZhbHVlOiBzdHJpbmcpID0+IHtcbiAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc2V0QXR0cmlidXRlKCdkYXRhLWNvbG9yLXNjaGVtZScsIHZhbHVlKTtcbn07XG5jb25zdCBzZXRDb2xvclNjaGVtZU1ldGEgPSAodmFsdWU6IHN0cmluZykgPT4ge1xuICAgIGNvbnN0IG1ldGE6IEhUTUxNZXRhRWxlbWVudCA9IGRvY3VtZW50LmhlYWQucXVlcnlTZWxlY3RvcignbWV0YVtuYW1lPVwiY29sb3Itc2NoZW1lXCJdJyk7XG4gICAgaWYgKG1ldGEpIHtcbiAgICAgICAgbWV0YS5jb250ZW50ID0gdmFsdWU7XG4gICAgfVxufTtcblxuLy8gQnkgZGVmYXVsdCBpdCBzaG91bGQgYmUgZGFyayB0aGUgc2l0ZS4gU28gaWYgc2FpZCBicm93c2VyXG4vLyBkb2Vzbid0IGhhdmUgdGhpcyBtZWRpYSBxdWVyeSBpdCB3aWxsIG1hdGNoZXMgdG8gZmFsc2UgYW5kXG4vLyBzZXQgdGhlIHNpdGUgdG8gYSBkYXJrIGNvbG9yLXNjaGVtZS5cbmNvbnN0IGxpZ2h0U2NoZW1lID0gbWF0Y2hNZWRpYSgnKHByZWZlcnMtY29sb3Itc2NoZW1lOiBsaWdodCknKTtcbmNvbnN0IGhhbmRsZUNvbG9yU2NoZW1lID0gKCkgPT4ge1xuICAgIGlmIChsaWdodFNjaGVtZS5tYXRjaGVzKSB7XG4gICAgICAgIHNldENvbG9yU2NoZW1lQXR0cmlidXRlKCdsaWdodCcpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHNldENvbG9yU2NoZW1lQXR0cmlidXRlKCdkYXJrJyk7XG4gICAgfVxufTtcblxuZXhwb3J0IGZ1bmN0aW9uIEluaXRhbGl6ZUNvbG9yU2NoZW1lKGNvbG9yU2NoZW1lOiBVc2VyU2V0dGluZ3NbJ2NvbG9yU2NoZW1lJ10pIHtcbiAgICBzd2l0Y2ggKGNvbG9yU2NoZW1lKSB7XG4gICAgICAgIGNhc2UgJ2ZvbGxvdy1zeXN0ZW0nOiB7XG4gICAgICAgICAgICBoYW5kbGVDb2xvclNjaGVtZSgpO1xuICAgICAgICAgICAgc2V0Q29sb3JTY2hlbWVNZXRhKCdkYXJrIGxpZ2h0Jyk7XG4gICAgICAgICAgICAvLyBBcyBpdCBmb2xsb3dzIHRoZSBzeXN0ZW0gd2Ugc2hvdWxkIGxpc3RlbiBmb3IgYW55IGNoYW5nZXMuXG4gICAgICAgICAgICBpZiAoaXNNYXRjaE1lZGlhQ2hhbmdlRXZlbnRMaXN0ZW5lclN1cHBvcnRlZCkge1xuICAgICAgICAgICAgICAgIGxpZ2h0U2NoZW1lLmFkZEV2ZW50TGlzdGVuZXIoJ2NoYW5nZScsIGhhbmRsZUNvbG9yU2NoZW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgbGlnaHRTY2hlbWUuYWRkTGlzdGVuZXIoaGFuZGxlQ29sb3JTY2hlbWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cbiAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgIHNldENvbG9yU2NoZW1lQXR0cmlidXRlKGNvbG9yU2NoZW1lKTtcbiAgICAgICAgICAgIHNldENvbG9yU2NoZW1lTWV0YShjb2xvclNjaGVtZSk7XG4gICAgICAgICAgICAvLyBNYWtlIHN1cmUgdG8gcmVtb3ZlIHRoZSBldmVudCBsaXN0ZW5lci5cbiAgICAgICAgICAgIGlmIChpc01hdGNoTWVkaWFDaGFuZ2VFdmVudExpc3RlbmVyU3VwcG9ydGVkKSB7XG4gICAgICAgICAgICAgICAgbGlnaHRTY2hlbWUucmVtb3ZlRXZlbnRMaXN0ZW5lcignY2hhbmdlJywgaGFuZGxlQ29sb3JTY2hlbWUpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICBsaWdodFNjaGVtZS5yZW1vdmVMaXN0ZW5lcihoYW5kbGVDb2xvclNjaGVtZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBicmVhaztcbiAgICB9XG59XG4iLCAiZXhwb3J0IGZ1bmN0aW9uIFNoYXJlQnV0dG9uKCkge1xuICAgIGNvbnN0IGkgPSBkb2N1bWVudC5nZXRFbGVtZW50QnlJZCgnc2hhcmUnKSBhcyBIVE1MSW5wdXRFbGVtZW50O1xuICAgIGNvbnN0IHNoYXJlQnV0dG9uID0gZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ2J0bi1zaGFyZScpIGFzIEhUTUxCdXR0b25FbGVtZW50O1xuICAgIHNoYXJlQnV0dG9uICYmIHNoYXJlQnV0dG9uLmFkZEV2ZW50TGlzdGVuZXIoJ2NsaWNrJywgKCkgPT4ge1xuICAgICAgICBpLnNlbGVjdCgpO1xuICAgICAgICBkb2N1bWVudC5leGVjQ29tbWFuZCgnY29weScpO1xuICAgICAgICBpLmJsdXIoKTtcbiAgICAgICAgc2hhcmVCdXR0b24uY2xhc3NMaXN0LmFkZCgnY29waWVkJyk7XG4gICAgfSk7XG59XG4iLCAiY29uc3QgcmVhZHlTdGF0ZUxpc3RlbmVycyA9IG5ldyBTZXQ8KCkgPT4gdm9pZD4oKTtcblxuZXhwb3J0IGNvbnN0IGlzRE9NUmVhZHkgPSAoKSA9PiAnY29tcGxldGUnID09PSBkb2N1bWVudC5yZWFkeVN0YXRlIHx8ICdpbnRlcmFjdGl2ZScgPT09IGRvY3VtZW50LnJlYWR5U3RhdGU7XG5leHBvcnQgY29uc3QgYWRkRE9NUmVhZHlMaXN0ZW5lciA9IChsaXN0ZW5lcjogKCkgPT4gdm9pZCkgPT4gcmVhZHlTdGF0ZUxpc3RlbmVycy5hZGQobGlzdGVuZXIpO1xuXG5pZiAoIWlzRE9NUmVhZHkoKSkge1xuICAgIGNvbnN0IG9uUmVhZHlTdGF0ZUNoYW5nZSA9ICgpID0+IHtcbiAgICAgICAgaWYgKGlzRE9NUmVhZHkoKSkge1xuICAgICAgICAgICAgZG9jdW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcigncmVhZHlzdGF0ZWNoYW5nZScsIG9uUmVhZHlTdGF0ZUNoYW5nZSk7XG4gICAgICAgICAgICByZWFkeVN0YXRlTGlzdGVuZXJzLmZvckVhY2goKGxpc3RlbmVyKSA9PiBsaXN0ZW5lcigpKTtcbiAgICAgICAgICAgIHJlYWR5U3RhdGVMaXN0ZW5lcnMuY2xlYXIoKTtcbiAgICAgICAgfVxuICAgIH07XG4gICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcigncmVhZHlzdGF0ZWNoYW5nZScsIG9uUmVhZHlTdGF0ZUNoYW5nZSk7XG59XG5cbmV4cG9ydCBjb25zdCByZW1vdmVFbGVtZW50ID0gKGVsZW1lbnQ6IEhUTUxFbGVtZW50KSA9PiB7XG4gICAgZWxlbWVudCAmJiBlbGVtZW50LnJlbW92ZSgpO1xufTtcbiIsICJpbXBvcnQge3JlbW92ZUVsZW1lbnR9IGZyb20gJy4vdXRpbHMvZG9tJztcblxuY29uc3QgZmlsbEluZm9ybWF0aW9uT25Gb3JtID0gKGtleTogc3RyaW5nLCB2YWx1ZTogc3RyaW5nKSA9PiB7XG4gICAgaWYgKCF2YWx1ZSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGNvbnN0IGVsZW1lbnQ6IEhUTUxJbnB1dEVsZW1lbnQgPSBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3IoYGZvcm0gW25hbWU9XCIke2tleX1cImApO1xuICAgIGlmICghZWxlbWVudCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGVsZW1lbnQudmFsdWUgPSB2YWx1ZTtcbn07XG5cbmNvbnN0IERFRkFVTFRfVVNFUlNUWUxFX01FVEEgPSBgLyogPT1Vc2VyU3R5bGU9PVxuQG5hbWUgICAgICAgICAgIEEgbmV3IHVzZXJzdHlsZSFcbkBuYW1lc3BhY2UgICAgICB1c2Vyc3R5bGVzLndvcmxkXG5AdmVyc2lvbiAgICAgICAgMS4wLjBcbj09L1VzZXJTdHlsZT09ICovXFxuYDtcblxuY29uc3QgaGFuZGxlU291cmNlQ29kZSA9IChzb3VyY2VDb2RlOiBzdHJpbmcpID0+IHtcbiAgICBpZiAoIXNvdXJjZUNvZGUpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBpZiAoIS9cXC9cXCogKj89PVVzZXJTdHlsZT09L2cudGVzdChzb3VyY2VDb2RlKSkge1xuICAgICAgICBzb3VyY2VDb2RlID0gYCR7REVGQVVMVF9VU0VSU1RZTEVfTUVUQX0ke3NvdXJjZUNvZGV9YDtcbiAgICB9XG4gICAgZmlsbEluZm9ybWF0aW9uT25Gb3JtKCdjb2RlJywgc291cmNlQ29kZSk7XG59O1xuXG5jb25zdCBvbk1lc3NhZ2UgPSAoZXY6IE1lc3NhZ2VFdmVudDxhbnk+KSA9PiB7XG4gICAgaWYgKCFldi5kYXRhIHx8ICFldi5kYXRhLnR5cGUpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBjb25zdCB7dHlwZSwgZGF0YX0gPSBldi5kYXRhO1xuICAgIHN3aXRjaCAodHlwZSkge1xuICAgICAgICBjYXNlICd1c3ctcmVtb3ZlLXN0eWx1cy1idXR0b24nOiB7XG4gICAgICAgICAgICByZW1vdmVFbGVtZW50KGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoJ2Ejc3R5bHVzJykpO1xuICAgICAgICB9XG4gICAgICAgIGNhc2UgJ3Vzdy1maWxsLW5ldy1zdHlsZSc6IHtcbiAgICAgICAgICAgIGlmICgnL2FwaS9vYXV0aC9hdXRob3JpemVfc3R5bGUvbmV3JyAhPT0gd2luZG93LmxvY2F0aW9uLnBhdGhuYW1lIHx8ICFkYXRhKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZmlsbEluZm9ybWF0aW9uT25Gb3JtKCduYW1lJywgZGF0YVsnbmFtZSddKTtcbiAgICAgICAgICAgIGNvbnN0IG1ldGFEYXRhID0gZGF0YVsnbWV0YWRhdGEnXTtcbiAgICAgICAgICAgIGlmIChtZXRhRGF0YSkge1xuICAgICAgICAgICAgICAgIGZpbGxJbmZvcm1hdGlvbk9uRm9ybSgnZGVzY3JpcHRpb24nLCBtZXRhRGF0YVsnZGVzY3JpcHRpb24nXSk7XG4gICAgICAgICAgICAgICAgZmlsbEluZm9ybWF0aW9uT25Gb3JtKCdsaWNlbnNlJywgbWV0YURhdGFbJ2xpY2Vuc2UnXSk7XG4gICAgICAgICAgICAgICAgZmlsbEluZm9ybWF0aW9uT25Gb3JtKCdob21lcGFnZScsIG1ldGFEYXRhWydsaWNlbnNlJ10pO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgaGFuZGxlU291cmNlQ29kZShkYXRhWydzb3VyY2VDb2RlJ10pO1xuXG4gICAgICAgIH1cbiAgICB9XG59O1xuXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbWVzc2FnZScsIG9uTWVzc2FnZSk7XG5cbmV4cG9ydCBmdW5jdGlvbiBCcm9hZGNhc3RSZWFkeSgpIHtcbiAgICB3aW5kb3cuZGlzcGF0Y2hFdmVudChuZXcgTWVzc2FnZUV2ZW50KCdtZXNzYWdlJywge1xuICAgICAgICBkYXRhOiB7dHlwZTogJ3Vzdy1yZWFkeSd9LFxuICAgICAgICBvcmlnaW46ICdodHRwczovL3VzZXJzdHlsZXMud29ybGQnXG4gICAgfSkpO1xufVxuIiwgImV4cG9ydCBpbnRlcmZhY2UgVXNlclNldHRpbmdzIHtcbiAgICBjb2xvclNjaGVtZTogJ2RhcmsnIHwgJ2xpZ2h0JyB8ICdmb2xsb3ctc3lzdGVtJztcbn1cblxuY29uc3QgREVGQVVMVF9TRVRUSU5HUzogVXNlclNldHRpbmdzID0ge1xuICAgIGNvbG9yU2NoZW1lOiAnZm9sbG93LXN5c3RlbScsXG59O1xuXG5jb25zdCBsb2NhbFN0b3JhZ2VLZXkgPSAndXNlci1wcmVmZXJlbmNlcyc7XG5sZXQgc2V0dGluZ3M6IFVzZXJTZXR0aW5ncyA9IG51bGw7XG5cbmV4cG9ydCBmdW5jdGlvbiBnZXRTZXR0aW5ncygpOiBVc2VyU2V0dGluZ3Mge1xuICAgIGlmIChzZXR0aW5ncykge1xuICAgICAgICByZXR1cm4gc2V0dGluZ3M7XG4gICAgfVxuICAgIGNvbnN0IE1heWJlU2V0dGluZ3MgPSBsb2NhbFN0b3JhZ2UuZ2V0SXRlbShsb2NhbFN0b3JhZ2VLZXkpO1xuICAgIGlmICghTWF5YmVTZXR0aW5ncykge1xuICAgICAgICBsb2NhbFN0b3JhZ2Uuc2V0SXRlbShsb2NhbFN0b3JhZ2VLZXksIEpTT04uc3RyaW5naWZ5KERFRkFVTFRfU0VUVElOR1MpKTtcbiAgICAgICAgcmV0dXJuIERFRkFVTFRfU0VUVElOR1M7XG4gICAgfVxuICAgIGNvbnN0IHNhdmVkU2V0dGluZ3MgPSBnZXRWYWxpZGF0ZWRPYmplY3QoSlNPTi5wYXJzZShNYXliZVNldHRpbmdzKSwgREVGQVVMVF9TRVRUSU5HUyk7XG5cbiAgICAvLyBEYXRhIG1pZ3JhdGlvbiwganVzdCB0byBiZSBzdXJlIGlmIGFueSBuZXcgc2V0dGluZyBhcmUgYWRkZWQuXG4gICAgLy8gV2Ugc2hvdWxkIGluY2x1ZGUgdGhlIGRlZmF1bHQgdmFsdWUgYW5kIHNhdmUgaXQuXG4gICAgc2V0dGluZ3MgPSB7Li4uREVGQVVMVF9TRVRUSU5HUywgLi4uc2F2ZWRTZXR0aW5nc307XG4gICAgbG9jYWxTdG9yYWdlLnNldEl0ZW0obG9jYWxTdG9yYWdlS2V5LCBKU09OLnN0cmluZ2lmeShzZXR0aW5ncykpO1xuXG4gICAgcmV0dXJuIHNldHRpbmdzO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc3RvcmVOZXdTZXR0aW5ncyhuZXdTZXR0aW5nczogUGFydGlhbDxVc2VyU2V0dGluZ3M+KSB7XG4gICAgc2V0dGluZ3MgPSB7Li4uc2V0dGluZ3MsIC4uLm5ld1NldHRpbmdzfTtcbiAgICBsb2NhbFN0b3JhZ2Uuc2V0SXRlbShsb2NhbFN0b3JhZ2VLZXksIEpTT04uc3RyaW5naWZ5KHNldHRpbmdzKSk7XG59XG5cbi8vIEEgbmllY2UgZnVuY3Rpb24gdG8gbWFrZSBzdXJlIHRoYXQgdGhlIHJldHVybiBvYmplY3Qgd2lsbCBvbmx5IHJldHVyblxuLy8ga2V5J3MgdGhhdCB0aGUgY29tcGFyZSBvYmplY3QgaGF2ZS5cbmZ1bmN0aW9uIGdldFZhbGlkYXRlZE9iamVjdDxUPihzb3VyY2U6IGFueSwgY29tcGFyZTogVCk6IFBhcnRpYWw8VD4ge1xuICAgIGNvbnN0IHJlc3VsdCA9IHt9O1xuICAgIGlmIChudWxsID09IHNvdXJjZSB8fCAnb2JqZWN0JyAhPT0gdHlwZW9mIHNvdXJjZSB8fCBBcnJheS5pc0FycmF5KHNvdXJjZSkpIHtcbiAgICAgICAgcmV0dXJuIG51bGw7XG4gICAgfVxuICAgIE9iamVjdC5rZXlzKHNvdXJjZSkuZm9yRWFjaCgoa2V5KSA9PiB7XG4gICAgICAgIGNvbnN0IHZhbHVlID0gc291cmNlW2tleV07XG4gICAgICAgIGlmIChudWxsID09IHZhbHVlIHx8IG51bGwgPT0gY29tcGFyZVtrZXldKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgY29uc3QgYXJyYXkxID0gQXJyYXkuaXNBcnJheSh2YWx1ZSk7XG4gICAgICAgIGNvbnN0IGFycmF5MiA9IEFycmF5LmlzQXJyYXkoY29tcGFyZVtrZXldKTtcbiAgICAgICAgaWYgKGFycmF5MSB8fCBhcnJheTIpIHtcbiAgICAgICAgICAgIGlmIChhcnJheTEgJiYgYXJyYXkyKSB7XG4gICAgICAgICAgICAgICAgcmVzdWx0W2tleV0gPSB2YWx1ZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSBlbHNlIGlmICgnb2JqZWN0JyA9PT0gdHlwZW9mIHZhbHVlICYmICdvYmplY3QnID09PSB0eXBlb2YgY29tcGFyZVtrZXldKSB7XG4gICAgICAgICAgICByZXN1bHRba2V5XSA9IGdldFZhbGlkYXRlZE9iamVjdCh2YWx1ZSwgY29tcGFyZVtrZXldKTtcbiAgICAgICAgfSBlbHNlIGlmICh0eXBlb2YgdmFsdWUgPT09IHR5cGVvZiBjb21wYXJlW2tleV0pIHtcbiAgICAgICAgICAgIHJlc3VsdFtrZXldID0gdmFsdWU7XG4gICAgICAgIH1cbiAgICB9KTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuIiwgImltcG9ydCB0eXBlIHtVc2VyU2V0dGluZ3N9IGZyb20gJy4vdXRpbHMvc3RvcmFnZSc7XG5pbXBvcnQge3N0b3JlTmV3U2V0dGluZ3N9IGZyb20gJy4vdXRpbHMvc3RvcmFnZSc7XG5cbmNvbnN0IFBSRUZJWCA9ICd1c3Itc2V0dGluZ3MnO1xuXG5leHBvcnQgZnVuY3Rpb24gU2V0VmFsdWVzKHNldHRpbmdzOiBVc2VyU2V0dGluZ3MpIHtcbiAgICBpZiAoIXdpbmRvdy5sb2NhdGlvbi5wYXRobmFtZS5zdGFydHNXaXRoKCcvYWNjb3VudCcpKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG4gICAgKGRvY3VtZW50LmdldEVsZW1lbnRCeUlkKGAke1BSRUZJWH0tLWNvbG9yLXNjaGVtZWApIGFzIEhUTUxTZWxlY3RFbGVtZW50KS52YWx1ZSA9IHNldHRpbmdzLmNvbG9yU2NoZW1lO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2F2ZVVzZXJTZXR0aW5nc0J1dHRvbihvblNldHRpbmdzVXBkYXRlOiAoKSA9PiB2b2lkKSB7XG4gICAgY29uc3Qgc2F2ZUJ1dHRvbiA9IGRvY3VtZW50LmdldEVsZW1lbnRCeUlkKGAke1BSRUZJWH0tLXNhdmVgKSBhcyBIVE1MQnV0dG9uRWxlbWVudDtcbiAgICBzYXZlQnV0dG9uICYmIHNhdmVCdXR0b24uYWRkRXZlbnRMaXN0ZW5lcignY2xpY2snLCAoKSA9PiB7XG4gICAgICAgIGNvbnN0IG5ld1NldHRpbmdzOiBQYXJ0aWFsPFVzZXJTZXR0aW5ncz4gPSB7fTtcblxuICAgICAgICBuZXdTZXR0aW5ncy5jb2xvclNjaGVtZSA9XG4gICAgICAgICAgICAoZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoYCR7UFJFRklYfS0tY29sb3Itc2NoZW1lYCkgYXMgSFRNTFNlbGVjdEVsZW1lbnQpLnZhbHVlIGFzIGFueTtcblxuICAgICAgICBzdG9yZU5ld1NldHRpbmdzKG5ld1NldHRpbmdzKTtcbiAgICAgICAgb25TZXR0aW5nc1VwZGF0ZSgpO1xuICAgIH0pO1xufVxuIiwgImltcG9ydCB7SW5pdGFsaXplQ29sb3JTY2hlbWUgYXMgaW5pdGFsaXplQ29sb3JTY2hlbWV9IGZyb20gJy4vY29sb3Itc2NoZW1lJztcbmltcG9ydCB7U2hhcmVCdXR0b259IGZyb20gJy4vc2hhcmUtYnV0dG9uJztcbmltcG9ydCB7QnJvYWRjYXN0UmVhZHl9IGZyb20gJy4vdGhpcmQtcGFydHknO1xuaW1wb3J0IHtTYXZlVXNlclNldHRpbmdzQnV0dG9uLCBTZXRWYWx1ZXN9IGZyb20gJy4vdXNlci1zZXR0aW5ncyc7XG5pbXBvcnQge2FkZERPTVJlYWR5TGlzdGVuZXIsIGlzRE9NUmVhZHl9IGZyb20gJy4vdXRpbHMvZG9tJztcbmltcG9ydCB7Z2V0U2V0dGluZ3N9IGZyb20gJy4vdXRpbHMvc3RvcmFnZSc7XG5cbmNvbnN0IFdoZW5ET01SZWFkeSA9ICgpID0+IHtcbiAgICBTaGFyZUJ1dHRvbigpO1xuICAgIEJyb2FkY2FzdFJlYWR5KCk7XG4gICAgU2F2ZVVzZXJTZXR0aW5nc0J1dHRvbihvblNldHRpbmdzVXBkYXRlKTtcbiAgICBTZXRWYWx1ZXMoZ2V0U2V0dGluZ3MoKSk7XG59O1xuXG4vLyBXaGVuRE9NUmVhZHkgY29udGFpbnMgY29kZSB0aGF0IG9ubHkgc2hvdWxkIGJlIGhhbmRsZVxuLy8gd2hlbiB0aGUgRE9NIGlzIHJlYWR5IHRvIGdvLlxuLy8gQW55IG90aGVyIGNvZGUgc2hvdWxkbid0IGRlcGVuZCBvbiB0aGlzIHNldHVwIGZ1bmN0aW9uLlxuaWYgKGlzRE9NUmVhZHkoKSkge1xuICAgIFdoZW5ET01SZWFkeSgpO1xufSBlbHNlIHtcbiAgICBhZGRET01SZWFkeUxpc3RlbmVyKFdoZW5ET01SZWFkeSk7XG59XG5cbi8vIE9uY2Ugc2V0dGluZ3MgdXBkYXRlIHdlIHNob3VsZCByZWluc3RhbGl6ZSBhbnkgZnVuY3Rpb25hbGxpdHkuXG4vLyBUaGF0IHJlbGllcyBvbiB0aGlzIHNldHRpbmdzLlxuY29uc3Qgb25TZXR0aW5nc1VwZGF0ZSA9ICgpID0+IHtcbiAgICBjb25zdCBzZXR0aW5ncyA9IGdldFNldHRpbmdzKCk7XG4gICAgaW5pdGFsaXplQ29sb3JTY2hlbWUoc2V0dGluZ3MuY29sb3JTY2hlbWUpO1xufTtcblxuLy8gSW5pdGFsaXplIGZ1bmN0aW9ucyB0aGF0IHJlcXVpcmVzIHNldHRpbmdzIGFuZCBkb24ndCBkZXBlbmQgb24gdGhlIERPTS5cbi8vIE5vdGUgdGhhdCB3ZSBkb24ndCBzYXZlIGdldFNldHRpbmdzKCkgcmVzdWx0LCBhcyB0aGlzIGluaXRhbGl6ZSBpcyBhIDEgdGltZSB0aGluZ1xuLy8gQW5kIGhhdmluZyBpdCBzaXQgaW4gdGhlIG1lbW9yeSBpcyBraW5kYSB1c2VsZXNzLlxuaW5pdGFsaXplQ29sb3JTY2hlbWUoZ2V0U2V0dGluZ3MoKS5jb2xvclNjaGVtZSk7XG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFPLE1BQU0sMkNBQ1QsQUFBZSxPQUFPLG1CQUF0QixjQUNBLEFBQWUsT0FBTyxlQUFlLFVBQVUscUJBQS9DOzs7QUNDSixNQUFNLDBCQUEwQixDQUFDLFVBQWtCO0FBQy9DLGFBQVMsZ0JBQWdC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
