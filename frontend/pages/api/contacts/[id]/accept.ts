import type { NextApiRequest, NextApiResponse } from "next";
import { getServerSession } from "next-auth/next";
import { authOptions } from "../../auth/[...nextauth]";

let backend = process.env.BACKEND_URL as string;

const getHeaders = (id: string) => {
  return {
    "Content-Type": "application/json",
    "USER-ID": id,
    "BACKEND-KEY": process.env.BACKEND_KEY as string,
  };
};

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST") {
    return res.status(405).json({ error: "Method not allowed" });
  }

  const session = await getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: "Unauthorized" });
  }

  const { id } = req.query;

  const response = await fetch(backend + `v1/contacts/${id}/accept`, {
    method: "POST",
    headers: getHeaders(session.providerAccountId),
  });

  if (response.status !== 200) {
    const errorText = await response.text();
    return res.status(response.status).json({ error: errorText || "Failed to accept contact" });
  }

  const data = await response.json();
  return res.status(200).json(data);
}
