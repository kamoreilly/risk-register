import React, { createContext, useCallback, useContext, useEffect, useMemo } from "react";
import { Uniwind, useUniwind } from "uniwind";

type ThemeName = "light" | "dark";

type AppThemeContextType = {
  currentTheme: string;
  isLight: boolean;
  isDark: boolean;
  setTheme: (theme: ThemeName) => void;
  toggleTheme: () => void;
};

const AppThemeContext = createContext<AppThemeContextType | undefined>(undefined);

export const AppThemeProvider = ({ children }: { children: React.ReactNode }) => {
  const { theme } = useUniwind();

  useEffect(() => {
    if (theme !== "dark") {
      Uniwind.setTheme("dark");
    }
  }, [theme]);

  const isLight = useMemo(() => {
    return theme === "light";
  }, [theme]);

  const isDark = useMemo(() => {
    return theme === "dark";
  }, [theme]);

  const setTheme = useCallback((newTheme: ThemeName) => {
    // Theme is locked to dark.
    Uniwind.setTheme("dark");
  }, []);

  const toggleTheme = useCallback(() => {
    // Theme is locked to dark.
    Uniwind.setTheme("dark");
  }, [theme]);

  const value = useMemo(
    () => ({
      currentTheme: theme,
      isLight,
      isDark,
      setTheme,
      toggleTheme,
    }),
    [theme, isLight, isDark, setTheme, toggleTheme],
  );

  return <AppThemeContext.Provider value={value}>{children}</AppThemeContext.Provider>;
};

export function useAppTheme() {
  const context = useContext(AppThemeContext);
  if (!context) {
    throw new Error("useAppTheme must be used within AppThemeProvider");
  }
  return context;
}
