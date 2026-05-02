import p1 from "@/assets/p1.jpg";
import p2 from "@/assets/p2.jpg";
import p3 from "@/assets/p3.jpg";
import p4 from "@/assets/p4.jpg";

export type Category = {
  slug: string;
  name: string;
  description: string;
  image: string;
};

export type Product = {
  id: string;
  slug: string;
  name: string;
  price: number;
  image: string;
  categorySlug: string;
  tag?: string;
  description: string;
};

export const categories: Category[] = [
  { slug: "audio", name: "Audio", description: "Headphones, speakers, microphones.", image: p1 },
  { slug: "wearables", name: "Wearables", description: "Watches and trackers.", image: p2 },
  { slug: "home", name: "Smart Home", description: "Lighting, climate, security.", image: p3 },
  { slug: "accessories", name: "Accessories", description: "Cables, stands, sleeves.", image: p4 },
];

export const products: Product[] = [
  { id: "p_001", slug: "wave-headphones", name: "Wave Headphones", price: 189, image: p1, categorySlug: "audio", tag: "New", description: "Over-ear, 40h battery, adaptive ANC." },
  { id: "p_002", slug: "studio-mic", name: "Studio Microphone", price: 129, image: p2, categorySlug: "audio", description: "USB-C condenser with monitoring." },
  { id: "p_003", slug: "pulse-watch", name: "Pulse Watch", price: 249, image: p3, categorySlug: "wearables", tag: "Popular", description: "GPS, ECG, 7-day battery." },
  { id: "p_004", slug: "lumen-bulb", name: "Lumen Smart Bulb", price: 32, image: p4, categorySlug: "home", description: "16M colors, schedules, voice control." },
  { id: "p_005", slug: "atlas-speaker", name: "Atlas Speaker", price: 219, image: p1, categorySlug: "audio", description: "360° sound, room calibration." },
  { id: "p_006", slug: "trace-band", name: "Trace Band", price: 89, image: p2, categorySlug: "wearables", description: "Sleep, stress, recovery insights." },
  { id: "p_007", slug: "nest-thermostat", name: "Nest Thermostat", price: 179, image: p3, categorySlug: "home", tag: "Limited", description: "Learning thermostat, energy reports." },
  { id: "p_008", slug: "fold-stand", name: "Fold Laptop Stand", price: 49, image: p4, categorySlug: "accessories", description: "Aluminum, six height settings." },
];

export type OrderStatus = "processing" | "shipped" | "delivered" | "cancelled";

export type Order = {
  id: string;
  date: string;
  total: number;
  status: OrderStatus;
  items: { name: string; qty: number; price: number; image: string }[];
};

export const sampleOrders: Order[] = [
  {
    id: "ORD-10293",
    date: "2025-04-12",
    total: 318,
    status: "shipped",
    items: [
      { name: "Wave Headphones", qty: 1, price: 189, image: p1 },
      { name: "Lumen Smart Bulb", qty: 4, price: 32, image: p4 },
    ],
  },
  {
    id: "ORD-10241",
    date: "2025-03-28",
    total: 249,
    status: "delivered",
    items: [{ name: "Pulse Watch", qty: 1, price: 249, image: p3 }],
  },
  {
    id: "ORD-10198",
    date: "2025-02-14",
    total: 89,
    status: "delivered",
    items: [{ name: "Trace Band", qty: 1, price: 89, image: p2 }],
  },
];

export const eur = (n: number) => `€${n.toFixed(0)}`;
