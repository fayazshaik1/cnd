import React, { Component }  from 'react';
import { hot } from 'react-hot-loader';
import Logo from './Logo.jsx';

import tvData from './data/tv.json';
import movieData from './data/movie.json';

import cindyAvatar from './assets/images/cindy.jpg';
import narcosBackground from './assets/images/narcos-bg.jpg';
import narcosLogo from './assets/images/narcos-logo.png';

import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>
        <header className="Header">
          <Logo />
          <Navigation />
          <UserProfile />
        </header>
        <Hero />
        <TitleList title="Top TV picks for Cindy" content='tv' />
        <TitleList title="Trending now" content='movie' />
      </div>
    );
  }
}


class Navigation extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div id="navigation" className="Navigation">
        <nav>
          <ul>
            <li>Browse</li>
            <li>My list</li>
            <li>Top picks</li>
            <li>Recent</li>
          </ul>
        </nav>
      </div>
    );
  }
}


class UserProfile extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div className="UserProfile">
        <div className="User">
          <div className="name">Cindy Lopez</div>
          <div className="image"><img src={cindyAvatar} alt="profile" /></div>
        </div>
      </div>
    );
  }
}


class Hero extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div id="hero" className="Hero" style={{backgroundImage: `url(${narcosBackground})`}}>
        <div className="content">
          <img className="logo" src={narcosLogo} alt="narcos background" />
          <h2>Season 2 now available</h2>
          <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Doloremque id quam sapiente 
            unde voluptatum alias vero debitis, magnam quis quod.</p>
          <div className="button-wrapper">
            <HeroButton primary={true} text="Watch now" />
            <HeroButton primary={false} text="+ My list" />
          </div>
        </div>
        <div className="overlay"></div>
      </div>
    );
  }
}


class HeroButton extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <a href="#" className="Button" data-primary={this.props.primary}>{this.props.text}</a>
    );
  }
}


class TitleList extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      data: [], 
      mounted: false
    };
  }

  loadContent() {
    let data;
    if (this.props.content === 'tv') {
      data = tvData;
    } else {
      data = movieData;
    }
    this.setState({ data: data });
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.content !== this.props.content && nextProps.content !== '') {
      this.setState({ 
        mounted: true,
        content: nextProps.content
      }, () => {
        this.loadContent();
      });
    }
  }

  componentDidMount() {
    if (this.props.url !== ''){
      this.loadContent();
      this.setState({ mounted: true });
    }
  }

  render() {
    var titles = '';
    if (this.state.data.results) {
      titles = this.state.data.results.map(function(title, i) {
        if (i < 5) {
          var name = '';
          var backDrop = `http://image.tmdb.org/t/p/original${title.backdrop_path}`;
          if (!title.name) {
            name = title.original_title;
          } else {
            name = title.name;
          }
          return (
            <Item 
              key={title.id} 
              title={name} 
              score={title.vote_average} 
              overview={title.overview} 
              backdrop={backDrop} 
            />
          );  
        } else {
          return (
            <div key={title.id}></div>
          );
        }
      }); 
    } 
    return (
      <div ref="titlecategory" className="TitleList" data-loaded={this.state.mounted}>
        <div className="Title">
          <h1>{this.props.title}</h1>
          <div className="titles-wrapper">
            {titles}
          </div>
        </div>
      </div>
    );
  }
}


class Item extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div className="Item" style={{ backgroundImage: 'url(' + this.props.backdrop + ')' }} >
        <div className="overlay">
          <div className="title">{this.props.title}</div>
          <div className="rating">{this.props.score} / 10</div>
          <div className="plot">{this.props.overview}</div>
          <ListToggle />
        </div>
      </div>
    );
  }
}


class ListToggle extends Component {
  constructor(props) {
    super(props);
    this.state = { toggled: false };
  }

  handleClick() {
    if(this.state.toggled === true) {
      this.setState({ toggled: false });
    } else {
      this.setState({ toggled: true }); 
    }
  }

  render() {
    return (
      <div onClick={this.handleClick} data-toggled={this.state.toggled} className="ListToggle">
        <div>
          <i className="fa fa-fw fa-plus"></i>
          <i className="fa fa-fw fa-check"></i>
        </div>
      </div>
    );
  }
}

export default hot(module)(App);
