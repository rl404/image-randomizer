import React from 'react';
import { Button, Grid, InputAdornment, TextField } from '@material-ui/core';
import PersonIcon from '@material-ui/icons/Person';
import VpnKeyIcon from '@material-ui/icons/VpnKey';
import { login, register } from '../api';

const Form = (props) => {
  const [state, setState] = React.useState({
    username: '',
    password: '',
    loading: false,
    message: '',
  });

  const submitForm = (e) => {
    e.preventDefault();
    setState({ ...state, loading: true, message: 'logging in...' });
    loginForm();
  }

  const loginForm = async () => {
    const result = await login({ username: state.username, password: state.password });
    if (result.status === 200) {
      setState({ ...state, loading: true, message: 'registering..' });
      props.setToken(state.username, result.data);
    } else if (result.status === 404) {
      registerForm();
    } else {
      setState({ ...state, message: result.message, loading: false });
    }
  }

  const registerForm = async () => {
    const result = await register({ username: state.username, password: state.password });
    if (result.status === 201) {
      loginForm();
    } else {
      setState({ ...state, message: result.message, loading: false });
    }
  }

  return (
    <form onSubmit={submitForm}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <TextField
            label='Username'
            placeholder='username'
            required
            variant='outlined'
            size='small'
            value={state.username}
            onChange={(e) => setState({ ...state, username: e.target.value })}
            InputProps={{
              startAdornment: (
                <InputAdornment position='start'>
                  <PersonIcon />
                </InputAdornment>
              )
            }} />
        </Grid>
        <Grid item xs={12}>
          <TextField
            label='Password'
            placeholder='password'
            required
            variant='outlined'
            type='password'
            size='small'
            value={state.password}
            onChange={(e) => setState({ ...state, password: e.target.value })}
            InputProps={{
              startAdornment: (
                <InputAdornment position='start'>
                  <VpnKeyIcon />
                </InputAdornment>
              )
            }} />
        </Grid>
        <Grid item xs={12}>
          <Button variant='contained' type='submit' color='primary'>
            Submit
          </Button>
          <span style={{ marginLeft: 50 }}>
            {state.message}
          </span>
        </Grid>
      </Grid>
    </form>
  );
};

export default Form;