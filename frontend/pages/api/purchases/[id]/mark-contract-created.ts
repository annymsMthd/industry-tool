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

  if (req.method === "POST") {
    const { id } = req.query;

    // Forward the request body to the backend
    const requestBody = typeof req.body === 'string' ? req.body : JSON.stringify(req.body);

    const response = await fetch(backend + `v1/purchases/${id}/mark-contract-created`, {
      method: "POST",
      headers: getHeaders(session.providerAccountId),
      body: requestBody,
    });

    // Read the response as text first to debug any parsing issues
    const responseText = await response.text();

    if (response.status !== 200) {
      try {
        const error = JSON.parse(responseText);
        return res.status(response.status).json(error);
      } catch (e) {
        console.error('Failed to parse error response:', responseText);
        return res.status(response.status).json({ error: responseText });
      }
    }

    try {
      const data = JSON.parse(responseText);
      return res.status(200).json(data);
    } catch (e) {
      console.error('Failed to parse success response:', responseText);
      return res.status(500).json({ error: 'Invalid response from backend', details: responseText });
    }
  }

  return res.status(405).json({ error: "Method not allowed" });
}
