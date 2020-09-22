import React from 'react';
import { Button, Grid, IconButton, TextField, Tooltip, Typography } from '@material-ui/core';
import { getList, getUserURL, updateList } from '../api';
import VisibilityIcon from '@material-ui/icons/Visibility';
import ClearIcon from '@material-ui/icons/Clear';
import AddIcon from '@material-ui/icons/Add';
import SaveIcon from '@material-ui/icons/Save';

const List = (props) => {
  const [state, setState] = React.useState({
    username: props.data.username,
    token: props.data.token,
    data: [],
    loading: true,
    preview: '',
    message: '',
  });

  React.useEffect(() => {
    if (state.token !== '' && state.username !== '' && state.loading) {
      getData();
    }
  });

  var timeout = 0;
  const getData = async () => {
    const result = await getList(state.username, state.token)
    if (result.status === 200) {
      setState({ ...state, data: result.data, loading: false });
    } else {
      setState({ ...state, message: result.message, loading: false });

      clearTimeout(timeout)
      timeout = setTimeout(() => { setState({ ...state, message: '' }) }, 5000);
    }
  };

  const updateData = async () => {
    const result = await updateList(state.token, state.data)
    if (result.status === 200) {
      setState({ ...state, message: 'done' })

      clearTimeout(timeout)
      timeout = setTimeout(() => { setState({ ...state, message: '' }) }, 5000);
    } else {
      setState({ ...state, message: result.message, loading: false });

      clearTimeout(timeout)
      timeout = setTimeout(() => { setState({ ...state, message: '' }) }, 5000);
    }
  }

  const addRow = () => {
    setState({ ...state, data: [...state.data, ''] })
  }

  const deleteRow = (i) => {
    setState({ ...state, data: [...state.data.slice(0, i), ...state.data.slice(i + 1)] })
  }

  const updateRow = (i, str) => {
    var tmp = state.data
    tmp[i] = str
    setState({ ...state, data: tmp })
  }

  const previewImg = (i) => {
    setState({ ...state, preview: state.data[i] })
  }

  const submitForm = (e) => {
    e.preventDefault();
    setState({ ...state, message: 'saving...' });
    updateData();
  }

  if (state.loading) {
    return 'wait a bit...';
  }

  return (
    <form onSubmit={submitForm}>
      <Grid container spacing={1}>
        <Grid item xs={12}>
          <Typography variant='h6'>
            {state.username}'s Images (<a href={getUserURL(state.username)} target='_blank' rel='noopener noreferrer'>{getUserURL(state.username)}</a>)
          </Typography>
        </Grid>
        <Grid item xs={12} container spacing={1} direction="row-reverse">
          <Grid item md={4} xs={12}>
            <img src={state.preview} style={{ width: '100%', position: 'sticky', top: 20 }} alt='preview' />
          </Grid>
          <Grid item md={8} xs={12} container spacing={1}>
            {state.data.map((image, i) => {
              return <Row
                key={i}
                index={i}
                image={image}
                delete={deleteRow}
                preview={previewImg}
                update={updateRow}
              />
            })}
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <Button variant="outlined" color="primary" startIcon={<AddIcon />} onClick={addRow} style={{ marginRight: 20 }} >
            Add row
          </Button>
          <Button variant='contained' color="primary" type='submit' startIcon={<SaveIcon />}>
            Submit
          </Button>
          <span style={{ marginLeft: 20 }}>
            {state.message}
          </span>
        </Grid>
      </Grid>
    </form>
  )
};

export default List;

const Row = (props) => {
  const image = props.image
  const i = props.index

  return (
    <>
      <Grid item xs={10}>
        <TextField
          placeholder='image url...'
          size='small'
          variant='outlined'
          fullWidth
          value={image}
          onChange={(e) => props.update(i, e.target.value)} />
      </Grid>
      <Grid item xs={2}>
        <Tooltip title='preview' placement='left'>
          <IconButton onClick={() => props.preview(i)}>
            <VisibilityIcon fontSize='small' />
          </IconButton>
        </Tooltip>
        <Tooltip title='delete' placement='right'>
          <IconButton onClick={() => props.delete(i)}>
            <ClearIcon fontSize='small' style={{ color: 'red' }} />
          </IconButton>
        </Tooltip>
      </Grid>
    </>
  )
}