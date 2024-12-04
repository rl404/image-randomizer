import IconButton from '@mui/material/IconButton';
import Grid from '@mui/material/Grid2';
import GitHubIcon from '@mui/icons-material/GitHub';

const Footer = () => {
  return (
    <Grid container justifyContent="flex-end">
      <Grid>
        <IconButton href="https://github.com/rl404/image-randomizer" target="_blank">
          <GitHubIcon />
        </IconButton>
      </Grid>
    </Grid>
  );
};

export default Footer;
