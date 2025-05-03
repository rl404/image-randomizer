export type Data = {
  status: number;
  message: string;
  data: Token;
};

export type Token = {
  access_token: string;
  refresh_token: string;
};

export async function POST(request: Request) {
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/register`, {
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
