const API_BASE = "";

export type Product = {
  id: string;
  slug: string;
  name: string;
  description: string;
  price: number;
  categoryId: string;
  categorySlug: string;
  sellerId?: string;
  image: string;
  stock?: number;
};

export type Category = {
  id: string;
  slug: string;
  name: string;
};

export type Order = {
  id: string;
  userId: string;
  status: "pending" | "paid" | "shipped" | "delivered" | "cancelled";
  total: number;
  createdAt: string;
  updatedAt: string;
  items: OrderItem[];
};

export type OrderItem = {
  id: string;
  orderId: string;
  productId: string;
  sellerId: string;
  quantity: number;
  price: number;
};

export type CreateOrderInput = {
  items: { product_id: string; quantity: number }[];
};

export type PaymentStatus =
  | "pending"
  | "processing"
  | "succeeded"
  | "failed"
  | "canceled"
  | "refunded";

export type CreatePaymentInput = {
  order_id: string;
  user_id: string;
  amount: number;
  currency: string;
  provider: string;
  idempotency_key: string;
};

export type PaymentInitResult = {
  paymentId: string;
  paymentUrl?: string;
  status: PaymentStatus;
  existing: boolean;
};

export type UserProfile = {
  id: string;
  name: string;
  email: string;
  role?: string;
  createdAt?: string;
  updatedAt?: string;
};

type RawProduct = {
  id: string;
  slug?: string;
  name: string;
  description?: string;
  price?: number | string;
  category_id?: string;
  categoryId?: string;
  category_slug?: string;
  categorySlug?: string;
  seller_id?: string;
  sellerId?: string;
  image_url?: string;
  imageURL?: string;
  image?: string;
  stock?: number;
};

type RawCategory = {
  id?: string;
  slug?: string;
  name: string;
};

type RawOrderItem = {
  id: string;
  order_id?: string;
  orderId?: string;
  product_id?: string;
  productId?: string;
  seller_id?: string;
  sellerId?: string;
  quantity: number;
  price?: number | string;
};

type RawOrder = {
  id: string;
  user_id?: string;
  userId?: string;
  status: Order["status"];
  total?: number | string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  items?: RawOrderItem[];
};

type RawUserProfile = {
  id: string;
  name: string;
  email: string;
  role?: string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
};

type RawPaymentInitResult = {
  payment_id?: string;
  paymentId?: string;
  payment_url?: string;
  paymentUrl?: string;
  status: PaymentStatus;
  existing?: boolean;
};

const getToken = () => {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
};

const authHeaders = () => {
  const token = getToken();
  if (!token) throw new Error("Please sign in first");
  return { Authorization: `Bearer ${token}` };
};

const getApiErrorMessage = async (res: Response, fallback: string) => {
  if (res.status === 401) {
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }
    return "Please sign in to continue";
  }

  const body = await res.json().catch(() => null);
  return body?.error || fallback;
};

const toSlug = (value: string) =>
  value
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");

const normalizeProduct = (raw: RawProduct): Product => ({
  id: raw.id,
  slug: raw.slug || raw.id,
  name: raw.name,
  description: raw.description || "",
  price: Number(raw.price || 0),
  categoryId: raw.category_id || raw.categoryId || "",
  categorySlug: raw.category_slug || raw.categorySlug || raw.category_id || raw.categoryId || "",
  sellerId: raw.seller_id || raw.sellerId,
  image: raw.image_url || raw.imageURL || raw.image || "https://placehold.co/400x500?text=No+Image",
  stock: raw.stock,
});

const normalizeCategory = (raw: RawCategory): Category => ({
  id: raw.id || raw.slug || toSlug(raw.name || "category"),
  slug: raw.slug || raw.id || toSlug(raw.name || "category"),
  name: raw.name,
});

const normalizeOrder = (raw: RawOrder): Order => ({
  id: raw.id,
  userId: raw.user_id || raw.userId || "",
  status: raw.status,
  total: Number(raw.total || 0),
  createdAt: raw.created_at || raw.createdAt || "",
  updatedAt: raw.updated_at || raw.updatedAt || "",
  items: (raw.items || []).map((item: RawOrderItem) => ({
    id: item.id,
    orderId: item.order_id || item.orderId || "",
    productId: item.product_id || item.productId || "",
    sellerId: item.seller_id || item.sellerId || "",
    quantity: item.quantity,
    price: Number(item.price || 0),
  })),
});

const normalizeUserProfile = (raw: RawUserProfile): UserProfile => ({
  id: raw.id,
  name: raw.name,
  email: raw.email,
  role: raw.role,
  createdAt: raw.created_at || raw.createdAt,
  updatedAt: raw.updated_at || raw.updatedAt,
});

const normalizePaymentInitResult = (raw: RawPaymentInitResult): PaymentInitResult => ({
  paymentId: raw.payment_id || raw.paymentId || "",
  paymentUrl: raw.payment_url || raw.paymentUrl,
  status: raw.status,
  existing: Boolean(raw.existing),
});

export const api = {
  // Auth
  register: async (data: { name: string; email: string; password: string }) => {
    const res = await fetch(`${API_BASE}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Registration failed"));
    return res.json();
  },
  login: async (data: { email: string; password: string }) => {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Login failed"));
    const json = await res.json();
    if (typeof window !== "undefined") {
      localStorage.setItem("token", json.accessToken);
    }
    return json;
  },
  logout: () => {
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }
    return fetch(`${API_BASE}/auth/logout`, { method: "POST" });
  },

  // Products
  getProducts: async () => {
    const res = await fetch(`${API_BASE}/products/`);
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to fetch products"));
    const json = await res.json();
    return (json as RawProduct[]).map(normalizeProduct);
  },
  getProduct: async (id: string) => {
    const res = await fetch(`${API_BASE}/products/${encodeURIComponent(id)}`);
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to fetch product"));
    return normalizeProduct(await res.json());
  },

  // Categories
  getCategories: async () => {
    const res = await fetch(`${API_BASE}/categories`);
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to fetch categories"));
    const json = await res.json();
    return (json as RawCategory[]).map(normalizeCategory);
  },

  // Users
  getMe: async () => {
    const res = await fetch(`${API_BASE}/users/me`, {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to fetch user"));
    return normalizeUserProfile(await res.json());
  },
  updateMe: async (data: { name: string; email: string }) => {
    const res = await fetch(`${API_BASE}/users/me`, {
      method: "PUT",
      headers: { "Content-Type": "application/json", ...authHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to update profile"));
    return normalizeUserProfile(await res.json());
  },
  createOrder: async (data: CreateOrderInput) => {
    const res = await fetch(`${API_BASE}/orders/`, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...authHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to create order"));
    return normalizeOrder(await res.json());
  },
  createPayment: async (data: CreatePaymentInput) => {
    const res = await fetch(`${API_BASE}/payments/init`, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...authHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to initialize payment"));
    return normalizePaymentInitResult(await res.json());
  },
  getMyOrders: async () => {
    const res = await fetch(`${API_BASE}/orders/my`, {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to fetch orders"));
    const json = await res.json();
    return (json as RawOrder[]).map(normalizeOrder);
  },
  getMyPayments: async () => {
    const profile = await api.getMe();
    const res = await fetch(`${API_BASE}/payments/users/${profile.id}`, {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error(await getApiErrorMessage(res, "Failed to fetch payments"));
    const json = await res.json();
    return (json as any[]).map((p) => ({
      id: p.id,
      orderId: p.order_id,
      amount: p.amount,
      currency: p.currency,
      provider: p.provider,
      status: p.status,
      createdAt: p.created_at,
    }));
  },
};
