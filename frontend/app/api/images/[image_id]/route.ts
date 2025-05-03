export async function PATCH(request: Request, { params }: { params: Promise<{ image_id: string }> }) {
  const { image_id } = await params;
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/images/${image_id}`, {
    headers: {
      Authorization: request.headers.get('authorization') || '',
      'Content-Type': 'application/json',
    },
    method: 'PATCH',
    body: JSON.stringify(await request.json()),
  });
  const data = await resp.json();
  return Response.json(data, { status: resp.status });
}

export async function DELETE(request: Request, { params }: { params: Promise<{ image_id: string }> }) {
  const { image_id } = await params;
  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/images/${image_id}`, {
    headers: {
      Authorization: request.headers.get('authorization') || '',
      'Content-Type': 'application/json',
    },
    method: 'DELETE',
  });
  const data = await resp.json();
  return Response.json(data, { status: resp.status });
}
