import { Button, Grid, IconButton, InputAdornment, Link, Stack, TextField, Tooltip, Typography } from '@mui/material';
import { NextPage } from 'next';
import { useRouter } from 'next/router';
import * as React from 'react';
import { Image } from '../types/Types';
import { axios2 } from '../utils/axios';
import { getAccessToken, getUsername } from '../utils/storage';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import Visibility from '@mui/icons-material/Visibility';
import DeleteIcon from '@mui/icons-material/Delete';
import { LoadingButton } from '@mui/lab';

const List: NextPage = () => {
  const router = useRouter();

  var timeout: NodeJS.Timeout;

  const [username, setUsername] = React.useState<string>('');
  const [images, setImages] = React.useState<Array<Image>>([]);
  const [error, setError] = React.useState<string>('');
  const [loading, setLoading] = React.useState<boolean>(true);
  const [preview, setPreview] = React.useState<string>('');

  React.useEffect(() => {
    if (!getAccessToken()) router.push('/');

    setUsername(getUsername());

    axios2
      .get('/api/images')
      .then((resp) => {
        const data: Array<Image> = resp.data.data;
        setImages(data);
        setLoading(false);
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            router.push('/');
            return;
          }

          setError(error.response.data.message);
          setLoading(false);
          return;
        }

        setError(error.message);
        setLoading(false);
      });
  }, []);

  const randomImgURL = `${process.env.NEXT_PUBLIC_API_HOST}/user/${username}/image.jpg`;

  const [copied, setCopied] = React.useState<boolean>(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(randomImgURL);
    setCopied(true);

    clearTimeout(timeout);
    timeout = setTimeout(() => {
      setCopied(false);
    }, 1000);
  };

  const newRowRef = React.useRef<HTMLButtonElement | null>(null);

  const addNewRow = () => {
    setImages([...images, { id: 0, user_id: 0, image: '' }]);
    newRowRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <Grid container spacing={2}>
      {loading ? (
        <Grid item>loading...</Grid>
      ) : error ? (
        <Grid item>{error}</Grid>
      ) : (
        <>
          <Grid item xs={12} sm={10} md={11}>
            <Typography variant="h5">
              {`${username}'s images`} (
              <Link href={randomImgURL} target="_blank" rel="noopener noreferrer">
                {randomImgURL}
              </Link>
              ){' '}
              <Tooltip title={!copied ? 'copy link' : 'copied!'} placement="right" arrow>
                <IconButton size="small" onClick={copyToClipboard}>
                  <ContentCopyIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Typography>
          </Grid>
          <Grid item xs={12} sm={2} md={1} sx={{ textAlign: 'right' }}>
            <Button href="/logout" variant="outlined" fullWidth>
              Logout
            </Button>
          </Grid>
          <Grid item xs={12} sm={3}>
            Preview
            {preview !== '' && (
              <img src={preview} alt="invalid image url" width="100%" style={{ position: 'sticky', top: 20 }} />
            )}
          </Grid>
          <Grid item xs={12} sm={9} container spacing={2}>
            {images.map((i) => {
              return <ImageRow key={i.id} image={i} setPreview={setPreview} />;
            })}
            <Grid item xs={12}>
              <Button fullWidth variant="outlined" onClick={addNewRow} ref={newRowRef}>
                Add new image
              </Button>
            </Grid>
          </Grid>
        </>
      )}
    </Grid>
  );
};

export default List;

const ImageRow = ({ image, setPreview }: { image: Image; setPreview: (link: string) => void }) => {
  const router = useRouter();

  const [formState, setFormState] = React.useState({
    id: image.id,
    image: image.image,
    error: '',
    deleted: false,
    loading: false,
  });

  const handlePreview = () => {
    setPreview(imageState.image);
  };

  const handleDelete = () => {
    if (formState.id === 0) {
      setFormState({ ...formState, deleted: true });
      return;
    }

    setFormState({ ...formState, loading: true });

    axios2
      .delete(`/api/images/${formState.id}`)
      .then((resp) => {
        setFormState({ ...formState, deleted: true, loading: false });
        setImageState({ ...imageState, showButton: false });
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            router.push('/');
            return;
          }

          setFormState({ ...formState, error: error.response.data.message, loading: false });
          return;
        }

        setFormState({ ...formState, error: error.message, loading: false });
      });
  };

  const [imageState, setImageState] = React.useState({
    image: image.image,
    showButton: false,
  });

  const handleChangeImage = (e: React.ChangeEvent<HTMLInputElement>) => {
    setImageState({ image: e.target.value, showButton: e.target.value !== formState.image });
  };

  const handleCancel = () => {
    setImageState({ image: formState.image, showButton: false });
  };

  const handleSubmit = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();

    setFormState({ ...formState, loading: true });

    if (formState.id === 0) {
      onCreate();
    } else {
      onUpdate();
    }
  };

  const onUpdate = () => {
    axios2
      .patch(`/api/images/${formState.id}`, {
        image: imageState.image,
      })
      .then((resp) => {
        setFormState({ ...formState, image: imageState.image, loading: false });
        setImageState({ ...imageState, showButton: false });
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            router.push('/');
            return;
          }

          setFormState({ ...formState, error: error.response.data.message, loading: false });
          return;
        }

        setFormState({ ...formState, error: error.message, loading: false });
      });
  };

  const onCreate = () => {
    axios2
      .post(`/api/images`, {
        image: imageState.image,
      })
      .then((resp) => {
        const data: Image = resp.data.data;
        setFormState({ ...formState, id: data.id, image: imageState.image, loading: false });
        setImageState({ ...imageState, showButton: false });
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.status === 401) {
            router.push('/');
            return;
          }

          setFormState({ ...formState, error: error.response.data.message, loading: false });
          return;
        }

        setFormState({ ...formState, error: error.message, loading: false });
      });
  };

  if (formState.id === 0 && formState.deleted) {
    return <></>;
  }

  return (
    <Grid item xs={12} sx={{ opacity: formState.deleted ? 0.3 : 1 }}>
      <form>
        <Stack direction="row" spacing={2}>
          <Tooltip
            title={
              <>
                {imageState.image.includes('imgur') &&
                  'Looks like you are hosting your image on Imgur. Your image may gets rate limitted by Imgur. '}
                {imageState.image.includes('discordapp') &&
                  'Looks like you are hosting your image on Discord. Discord link is not permanent anymore and can be expired. '}
                Try to host it on other site such as{' '}
                <Link href="https://postimages.org/" target="_blank" rel="noopener noreferrer">
                  postimages
                </Link>
                ,{' '}
                <Link href="https://imgbb.com/" target="_blank" rel="noopener noreferrer">
                  imgbb
                </Link>
                , or{' '}
                <Link href="https://github.com/" target="_blank" rel="noopener noreferrer">
                  github
                </Link>
                .
              </>
            }
            placement="left"
            arrow
            open={imageState.image.includes('imgur') || imageState.image.includes('discordapp')}
          >
            <TextField
              placeholder="http://your.image.url.com"
              required
              fullWidth
              size="small"
              disabled={formState.loading || formState.deleted}
              value={imageState.image}
              onChange={handleChangeImage}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <Stack direction="row" spacing={1}>
                      <Tooltip title="show preview" placement="left" arrow>
                        <IconButton onClick={handlePreview} size="small" disabled={formState.deleted}>
                          <Visibility fontSize="small" />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="delete" placement="right" arrow>
                        <IconButton onClick={handleDelete} edge="end" size="small" disabled={formState.deleted}>
                          <DeleteIcon color="error" fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </Stack>
                  </InputAdornment>
                ),
              }}
            />
          </Tooltip>
          {imageState.showButton && (
            <>
              <LoadingButton onClick={handleSubmit} type="submit" loading={formState.loading}>
                Update
              </LoadingButton>
              <LoadingButton onClick={handleCancel} loading={formState.loading}>
                Cancel
              </LoadingButton>
            </>
          )}
        </Stack>
      </form>
    </Grid>
  );
};
