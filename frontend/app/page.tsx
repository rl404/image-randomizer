import FormLogin from '@/src/components/FormLogin';
import Grid from '@mui/material/Grid';

export default function Home() {
  return (
    <Grid container justifyContent="center" alignItems="center">
      <Grid size={{ xs: 12, sm: 6, md: 3 }}>
        <FormLogin />
      </Grid>
    </Grid>
  );
}
