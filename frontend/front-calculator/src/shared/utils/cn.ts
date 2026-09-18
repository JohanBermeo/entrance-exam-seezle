type ClassValue = string | false | null | undefined

/** Une clases CSS ignorando valores falsy. */
export function cn(...parts: ClassValue[]): string {
  return parts.filter(Boolean).join(' ')
}
