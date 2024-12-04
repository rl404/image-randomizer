import Grid from '@mui/material/Grid2';
import Typography from '@mui/material/Typography';
import React from 'react';

const Header = () => {
  return (
    <Grid container spacing={2} sx={{ marginTop: 3 }}>
      <Grid size={12}>
        <Typography variant="h3">Image Randomizer</Typography>
      </Grid>
      <Grid size={12}>
        <Typography variant="subtitle1">
          Randoming image from your chosen image list. Can be used for randoming website background image, for example.
        </Typography>
        <Typography variant="subtitle1">{`Just create a username and password and that's it.`}</Typography>
      </Grid>
    </Grid>
  );
};

export default Header;
