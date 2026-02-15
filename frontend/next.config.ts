import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/backend/:path*",
        destination: `${process.env.BACKEND_URL || "http://localhost:8081"}/:path*`,
      },
    ];
  },
  async headers() {
    return [
      {
        source: "/backend/:path*",
        headers: [
          {
            key: "BACKEND-KEY",
            value: process.env.BACKEND_KEY || "",
          },
        ],
      },
    ];
  },
  reactStrictMode: true,
};

export default nextConfig;
