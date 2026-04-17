import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function buildQueryString<T extends Record<string, unknown>>(
  params?: T
): string {
  if (!params) return '';

  const searchParams = new URLSearchParams();

  (Object.keys(params) as Array<keyof T>).forEach((key) => {
    const value = params[key];
    if (value !== undefined && value !== null && value !== '') {
      searchParams.set(key as string, String(value));
    }
  });

  const qs = searchParams.toString();
  return qs ? `?${qs}` : '';
}
