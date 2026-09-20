export async function GET(request: Request, { params }: { params: Promise<{ username: string }> }) {
  const { username } = await params;
  return Response.json(`${process.env.NEXT_PUBLIC_API_HOST}/user/${username}/image.jpg`, { status: 200 });
}
