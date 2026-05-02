export type CartLine = {
  id: string;
  qty: number;
};

const CART_KEY = "cart";

export const readCart = (): CartLine[] => {
  if (typeof window === "undefined") return [];
  try {
    const value = localStorage.getItem(CART_KEY);
    if (!value) return [];
    const parsed = JSON.parse(value);
    if (!Array.isArray(parsed)) return [];
    return parsed.filter((line) => typeof line.id === "string" && typeof line.qty === "number");
  } catch {
    return [];
  }
};

export const writeCart = (lines: CartLine[]) => {
  if (typeof window === "undefined") return;
  localStorage.setItem(CART_KEY, JSON.stringify(lines));
};

export const clearCart = () => writeCart([]);

export const addToCart = (id: string, qty = 1) => {
  const lines = readCart();
  const existing = lines.find((line) => line.id === id);
  if (existing) {
    existing.qty += qty;
  } else {
    lines.push({ id, qty });
  }
  writeCart(lines);
};
