// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next';
import { Token } from '../../types/Types';

type Data = {
  status: number;
  message: string;
  data: Token;
};

export default async function handler(req: NextApiRequest, res: NextApiResponse<Data>) {
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/register`, {
    headers: {
      'Content-Type': 'application/json',
    },
    method: req.method,
    body: JSON.stringify(req.body),
  });
  const data = await resp.json();
  res.status(data.status).json(data);
}
