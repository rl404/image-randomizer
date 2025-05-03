export type Data = {
  status: number;
  message: string;
  data: Image[];
};

export type Image = {
  id: number;
  user_id: number;
  image: string;
};

export async function GET(request: Request) {
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/images`, {
    headers: {
      Authorization: request.headers.get('authorization') || '',
      'Content-Type': 'application/json',
    },
  });
  const data = await resp.json();
  return Response.json(data, { status: resp.status });
}

export async function POST(request: Request) {
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/images`, {
    headers: {
      Authorization: request.headers.get('authorization') || '',
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(await request.json()),
  });
  const data = await resp.json();
  return Response.json(data, { status: resp.status });
}
