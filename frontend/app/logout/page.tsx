'use client';

import { deleteStorage } from '@/src/utils/storage';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function Logout() {
  const router = useRouter();

  useEffect(() => {
    deleteStorage();
    router.push('/');
  }, []);

  return <>Logging out...</>;
}
