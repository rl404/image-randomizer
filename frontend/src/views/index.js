import React from 'react';
import Form from './Form';
import List from './List';

const Home = (props) => {
  const [state, setState] = React.useState({
    username: '',
    token: '',
  });
  const setToken = (username, token) => {
    setState({ username: username, token: token });
  };

  if (state.token === '' || state.username === '') {
    return <Form setToken={setToken} />
  } else {
    return <List data={state} />
  }
};

export default Home;