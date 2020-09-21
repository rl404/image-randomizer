import React from 'react';
import { Container, CssBaseline } from '@material-ui/core';
import {
  Redirect,
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';
import Header from './components/header';
import Footer from './components/footer';
import BackdropLoading from './components/loading/Backdrop';

const Home = React.lazy(() => import('./views'));

class App extends React.Component {
  render() {
    return (
      <main>
        <CssBaseline />
        <Container>
          <Header />
          <Router>
            <React.Suspense fallback={<BackdropLoading />}>
              <Switch>
                <Route path="/" name="Home" exact>
                  <Home />
                </Route>
                <Redirect from='/' to='/' />
              </Switch>
            </React.Suspense>
          </Router>
          <Footer />
        </Container>
      </main>
    );
  };
}

export default App;
