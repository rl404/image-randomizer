import { Grid, Typography } from '@mui/material';
import React from 'react';

const Header = () => {
  return (
    <Grid container spacing={2} sx={{ marginTop: 3 }}>
      <Grid item xs={12}>
        <Typography variant="h3">Image Randomizer</Typography>
      </Grid>
      <Grid item xs={12}>
        <Typography variant="subtitle1">
          Randoming image from your chosen image list. Can be used for randoming website background image, for example.
        </Typography>
        <Typography variant="subtitle1">{`Just create a username and password and that's it.`}</Typography>
      </Grid>
    </Grid>
  );
};

export default Header;
