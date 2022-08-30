// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next';
import { Image } from '../../../types/Types';

type Data = {
  status: number;
  message: string;
  data: Array<Image> | Image;
};

export default async function handler(req: NextApiRequest, res: NextApiResponse<Data>) {
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/images`, {
    headers: {
      Authorization: req.headers.authorization || '',
      'Content-Type': 'application/json',
    },
    method: req.method,
    body: req.method !== 'GET' ? JSON.stringify(req.body) : null,
  });
  const data = await resp.json();
  res.status(data.status).json(data);
}
