import React from 'react';
import { Divider, Grid, IconButton } from '@material-ui/core';
import GitHubIcon from '@material-ui/icons/GitHub';

const Footer = () => {
  return (
    <div style={{ marginTop: 50 }}>
      <Divider />
      <Grid container spacing={1} justify="flex-end">
        <Grid item>
          <a href='https://github.com/rl404/image-randomizer' target='_blank' rel='noopener noreferrer'>
            <IconButton>
              <GitHubIcon fontSize="small" />
            </IconButton>
          </a>
        </Grid>
      </Grid>
    </div>
  );
};

export default Footer;