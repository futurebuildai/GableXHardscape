import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

/**
 * Format an int64-cents value from the ERP API as a locale-formatted currency string.
 * ERP endpoints (orders, invoices) return money as integer cents; dividing by 100
 * and locale-formatting avoids the $7,388 → $73.88 bug.
 */
export function formatCents(cents: number): string {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
    }).format((cents || 0) / 100);
}
