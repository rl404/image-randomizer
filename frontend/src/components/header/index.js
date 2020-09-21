import React from 'react';
import { Divider, Typography } from '@material-ui/core';

const Header = () => {
  return (
    <div style={{ marginBottom: 50 }}>
      <Typography variant='h3'>
        Image Randomizer
      </Typography>
      <Typography variant='subtitle1'>
        Randoming image from your chosen image list.
        Can be used for randoming website background image, for example.
      </Typography>
      <Typography variant='subtitle1'>
        Just create a username and password and that's it.
      </Typography>
      <Divider />
    </div>
  );
};

export default Header;