import Grid from '@mui/material/Grid2';
import type { NextPage } from 'next';
import Head from 'next/head';
import { useRouter } from 'next/router';
import * as React from 'react';
import FormLogin from '../components/Form/FormLogin';
import { getAccessToken } from '../utils/storage';

const Home: NextPage = () => {
  const router = useRouter();

  React.useEffect(() => {
    if (getAccessToken()) router.push('/list');
  }, []);

  const title = 'Image Randomizer';
  const desc = 'Generate image randomly from your chosen image list';

  return (
    <>
      <Head>
        <title>{title}</title>
        <meta charSet="utf-8" />
        <meta name="title" content={title} />
        <meta name="description" content={desc} />
        <meta property="og:title" content={title} />
        <meta property="og:description" content={desc} />
        <meta property="og:type" content="website" />
        <meta property="og:site_name" content={title} />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="msapplication-TileColor" content="#3f51b5" />
        <meta name="theme-color" content="#ffffff" />
        <link rel="icon" href="/favicon.ico" />
        <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
        <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
        <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
        <link rel="manifest" href="/site.webmanifest" />
        <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#3f51b5" />
      </Head>
      <Grid container justifyContent="center" alignItems="center">
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <FormLogin />
        </Grid>
      </Grid>
    </>
  );
};

export default Home;
