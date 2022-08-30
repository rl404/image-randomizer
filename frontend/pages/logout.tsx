import { NextPage } from 'next';
import { useRouter } from 'next/router';
import * as React from 'react';
import { deleteStorage } from '../utils/storage';

const Logout: NextPage = () => {
  const router = useRouter();

  React.useEffect(() => {
    deleteStorage();
    router.push('/');
  }, []);

  return <>Logging out...</>;
};

export default Logout;
