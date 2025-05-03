import GitHubIcon from '@mui/icons-material/GitHub';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';

export default function Footer() {
  return (
    <Grid container justifyContent="flex-end">
      <Grid>
        <IconButton href="https://github.com/rl404/image-randomizer" target="_blank">
          <GitHubIcon />
        </IconButton>
      </Grid>
    </Grid>
  );
}
