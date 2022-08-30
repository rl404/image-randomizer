import { Grid, IconButton } from '@mui/material';
import GitHubIcon from '@mui/icons-material/GitHub';

const Footer = () => {
  return (
    <Grid container justifyContent="flex-end">
      <Grid item>
        <IconButton href="https://github.com/rl404/image-randomizer" target="_blank">
          <GitHubIcon />
        </IconButton>
      </Grid>
    </Grid>
  );
};

export default Footer;
