import KeyIcon from '@mui/icons-material/Key';
import PersonIcon from '@mui/icons-material/Person';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import LoadingButton from '@mui/lab/LoadingButton';
import Grid from '@mui/material/Grid2';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import axios from 'axios';
import { useRouter } from 'next/router';
import * as React from 'react';
import { Token } from '../../types/Types';
import { saveAccessToken, saveRefreshToken, saveUsername } from '../../utils/storage';

const FormLogin = () => {
  const router = useRouter();

  const [formState, setFormState] = React.useState({
    username: '',
    password: '',
    showPassword: false,
    error: '',
    loading: false,
  });

  const onChangeField = (name: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormState({ ...formState, [name]: e.target.value });
  };

  const togglePassword = () => {
    setFormState({ ...formState, showPassword: !formState.showPassword });
  };

  const onSubmit = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();

    if (formState.username === '') {
      setFormState({ ...formState, error: 'required field username', loading: false });
      return;
    }

    if (formState.password === '') {
      setFormState({ ...formState, error: 'required field password', loading: false });
      return;
    }

    setFormState({ ...formState, loading: true });
    onLogin();
  };

  const onLogin = async () => {
    await axios
      .post('/api/login', {
        username: formState.username,
        password: formState.password,
      })
      .then((resp) => {
        const data: Token = resp.data.data;
        saveAccessToken(data.access_token);
        saveRefreshToken(data.refresh_token);
        saveUsername(formState.username);
        router.push('/list');
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 404) {
            onRegister();
            return;
          }

          setFormState({ ...formState, error: error.response.data.message, loading: false });
          return;
        }

        setFormState({ ...formState, error: error.message, loading: false });
      });
  };

  const onRegister = async () => {
    await axios
      .post('/api/register', {
        username: formState.username,
        password: formState.password,
      })
      .then((resp) => {
        const data: Token = resp.data.data;
        saveAccessToken(data.access_token);
        saveRefreshToken(data.refresh_token);
        saveUsername(formState.username);
        router.push('/list');
      })
      .catch((error) => {
        if (error.response) {
          setFormState({ ...formState, error: error.response.data.message, loading: false });
          return;
        }

        setFormState({ ...formState, error: error.message, loading: false });
      });
  };

  return (
    <form>
      <Grid container spacing={2}>
        <Grid size={12}>
          <TextField
            label="Username"
            placeholder="username"
            required
            fullWidth
            variant="outlined"
            size="small"
            value={formState.username}
            onChange={onChangeField('username')}
            sx={{
              '& input:-webkit-autofill': {
                boxShadow: '0 0 0 1000px #121212 inset !important',
                WebkitTextFillColor: 'white !important',
                borderRadius: 0,
              },
            }}
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    <PersonIcon />
                  </InputAdornment>
                ),
              },
            }}
          />
        </Grid>
        <Grid size={12}>
          <TextField
            label="Password"
            placeholder="password"
            required
            fullWidth
            type={formState.showPassword ? 'text' : 'password'}
            variant="outlined"
            size="small"
            value={formState.password}
            onChange={onChangeField('password')}
            sx={{
              '& input:-webkit-autofill': {
                boxShadow: '0 0 0 1000px #121212 inset',
                WebkitTextFillColor: 'white',
                borderRadius: 0,
              },
            }}
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    <KeyIcon />
                  </InputAdornment>
                ),
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={togglePassword} edge="end">
                      {formState.showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              },
            }}
          />
        </Grid>
        <Grid size={12}>
          <LoadingButton variant="contained" type="submit" fullWidth onClick={onSubmit} loading={formState.loading}>
            Login or Register
          </LoadingButton>
        </Grid>
        {formState.error !== '' && (
          <Grid size={12}>
            <Typography color="error">{formState.error}</Typography>
          </Grid>
        )}
      </Grid>
    </form>
  );
};

export default FormLogin;
