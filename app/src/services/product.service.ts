import type { Product } from '../types/product';
import { fetchWithAuth } from './fetchClient';

const API_URL = import.meta.env.VITE_API_URL || '';

/**
 * Backend list endpoints return a paginated envelope:
 *   { data: Product[], total, limit, offset }
 * The service unwraps `.data` so callers can stay typed against Product[].
 * A defensive Array.isArray fallback covers any endpoint that returns the
 * array bare (legacy) or any future shape change.
 */
function unwrapList<T>(raw: unknown): T[] {
    if (Array.isArray(raw)) return raw as T[];
    if (raw && typeof raw === 'object' && Array.isArray((raw as { data?: unknown }).data)) {
        return (raw as { data: T[] }).data;
    }
    return [];
}

export const ProductService = {
    async getProducts(): Promise<Product[]> {
        const response = await fetchWithAuth(`${API_URL}/api/v1/products`);
        if (!response.ok) {
            throw new Error('Failed to fetch products');
        }
        return unwrapList<Product>(await response.json());
    },

    async createProduct(product: Omit<Product, 'id' | 'created_at' | 'updated_at'>): Promise<Product> {
        const response = await fetchWithAuth(`${API_URL}/api/v1/products`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to create product');
        }

        return response.json();
    },
};
