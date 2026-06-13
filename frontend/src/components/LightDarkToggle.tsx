import Sun from "/src/assets/sun.svg?react";
import Moon from "/src/assets/moon.svg?react";
import { Switch } from "./ui/switch";
import { useTheme } from "./ThemeProvider";

interface LightDarkToggleProps {}

export default function LightDarkToggle({}: LightDarkToggleProps) {
  const { theme, toggleTheme } = useTheme();
  return (
    <div className="flex items-center gap-2">
      <Sun className="size-5" />
      <Switch checked={theme === "dark"} onCheckedChange={toggleTheme} />
      <Moon className="size-5" />
    </div>
  );
}
