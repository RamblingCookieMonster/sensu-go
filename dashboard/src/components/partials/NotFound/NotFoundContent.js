import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import Grid from "@material-ui/core/Grid";
import Hidden from "@material-ui/core/Hidden";

const Frame = withStyles(() => ({
  root: {
    flexGrow: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
}))(({ classes, ...props }) => <div className={classes.root} {...props} />);

const GridContainer = withStyles(() => ({
  container: {
    maxWidth: 1080,
    justifyContent: "center",
  },
}))(props => <Grid container spacing={40} {...props} />);

const GridItem = withStyles(theme => ({
  item: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    color: theme.palette.text.hint,
  },
}))(props => <Grid item {...props} />);

const ContentContainer = withStyles(theme => ({
  root: {
    textAlign: "center",
    [theme.breakpoints.up("md")]: {
      textAlign: "left",
      paddingRight: theme.spacing.unit * 6,
    },
  },
}))(({ classes, ...props }) => <div className={classes.root} {...props} />);

const Title = withStyles(theme => ({
  root: {
    fontSize: 128,
    fontWeight: 100,
    color: theme.palette.text.secondary,
  },
}))(props => <Typography variant="title" {...props} />);

const Subheading = withStyles(theme => ({
  root: {
    fontSize: 24,
    fontWeight: "normal",
    color: theme.palette.text.hint,
  },
}))(props => <Typography variant="subheading" {...props} />);

const GraphicContainer = withStyles({
  root: {
    fontSize: 96,
  },
})(props => <Title {...props} />);

const styles = theme => ({
  body: {
    color: theme.palette.text.hint,
    "& a": {
      textDecoration: "underline",
    },
    "& a:visited": {
      color: "inherit",
    },
    "& a:link": {
      color: "inherit",
    },
  },
});

const symbols = ["✄···", "🤯", "🐳", "🚀", "🏗", "🚧"];

class NotFoundContent extends React.Component {
  static propTypes = {
    classes: PropTypes.object.isRequired,
    graphic: PropTypes.node,
    title: PropTypes.node,
    subtitle: PropTypes.node,
    content: PropTypes.node,
  };

  static defaultProps = {
    graphic: null,
    content: null,
    title: "404",
    subtitle: "The page you requested isn’t here.",
  };

  constructor(props) {
    super(props);
    this.symbol = symbols[Math.floor(Math.random() * symbols.length)];
  }

  renderDefaultContent = () => (
    <React.Fragment>
      <p>
        If you opened a link, it’s possible that the resource was deleted or you
        no longer have access.
      </p>
      <p>
        <a href="#back" onClick={() => window.history.back()}>
          Go back
        </a>{" "}
        or <a href="/">return home</a>.
      </p>
    </React.Fragment>
  );

  renderSymbol = () => (
    <span role="img" aria-label="not-found">
      {this.symbol}
    </span>
  );

  render() {
    return (
      <Frame>
        <GridContainer>
          <Hidden smDown>
            <GridItem sm={6}>
              <GraphicContainer>
                {this.props.graphic || this.symbol}
              </GraphicContainer>
            </GridItem>
          </Hidden>
          <GridItem xs={12} sm={6}>
            <ContentContainer>
              <Title>{this.props.title}</Title>
              <Subheading>{this.props.subtitle}</Subheading>
              <Typography
                component="div"
                variant="body1"
                color="inherit"
                className={this.props.classes.body}
              >
                {this.props.content || this.renderDefaultContent()}
              </Typography>
            </ContentContainer>
          </GridItem>
        </GridContainer>
      </Frame>
    );
  }
}

export default withStyles(styles)(NotFoundContent);
