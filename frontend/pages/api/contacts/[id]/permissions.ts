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
  const session = await getServerSession(req, res, authOptions);
  if (!session) {
    return res.status(401).json({ error: "Unauthorized" });
  }

  const { id } = req.query;

  if (req.method === "GET") {
    // Get permissions for contact
    const response = await fetch(backend + `v1/contacts/${id}/permissions`, {
      method: "GET",
      headers: getHeaders(session.providerAccountId),
    });

    if (response.status !== 200) {
      const errorText = await response.text();
      return res.status(response.status).json({ error: errorText || "Failed to get permissions" });
    }

    const data = await response.json();
    return res.status(200).json(data);
  } else if (req.method === "POST") {
    // Update permission
    const response = await fetch(backend + `v1/contacts/${id}/permissions`, {
      method: "POST",
      headers: getHeaders(session.providerAccountId),
      body: JSON.stringify(req.body),
    });

    if (response.status !== 200) {
      const errorText = await response.text();
      return res.status(response.status).json({ error: errorText || "Failed to update permission" });
    }

    return res.status(200).json({ success: true });
  } else {
    return res.status(405).json({ error: "Method not allowed" });
  }
}
