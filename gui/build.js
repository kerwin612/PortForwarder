const { exec, execSync } = require( 'child_process' );
const http = require( 'http' );
const os = require( 'os' );
const fs = require( 'fs' );

function isUrlReachable( url ) {
    return new Promise( ( resolve ) => {
        http.request( new URL( url ), { method: 'HEAD' }, ( res ) => {
            if ( res.statusCode >= 200 && res.statusCode < 300 ) {
                resolve( true );
            } else {
                resolve( false );
            }
        } ).on( 'error', () => {
            resolve( false );
        } ).end();
    } );
}

function checkStatik() {
  return new Promise( ( resolve, _reject ) => {
    const command = os.platform() === 'win32' ? 'where' : 'which';
    exec( `${command} statik`, ( error, _stdout, _stderr ) => {
      if ( error ) {
        resolve( false );
      } else {
        resolve( true );
      }
    } );
  } );
}

function installStatik() {
  try {
    execSync( 'go install github.com/rakyll/statik', { stdio: 'ignore' } );
    console.log( 'statik has been successfully installed.' );
    return true;
  } catch ( error ) {
    console.error( 'An error occurred during installing statik:', error );
  }
  return false;
}

( async function main() {
    if ( !( await checkStatik() ) && !( installStatik() ) ) {
      console.error( 'An exception occurred while checking or installing statik.' );
      process.exit( 1 );
    }
    const args = process.argv.slice( 2 );
    const action = args.length < 1 ? 'app' : args[0];
    switch ( action ) {
        case 'app':
            if ( !( fs.existsSync( '_statik' ) ) ) {
                if ( !( fs.existsSync( 'build' ) ) ) {
                    execSync( 'npm run build:gui', { stdio: 'inherit' } );
                }
                execSync( 'statik -f -src=build -p _statik -dest ./', { stdio: 'inherit' } );
            }
            execSync( 'go run ./main.go -ws=./_dev -ao=false', { stdio: 'inherit' } );
            break;
        case 'gui':
            let url = ( fs.existsSync( './_dev/pid' ) ) ? fs.readFileSync( './_dev/pid', 'utf8' ) : null;
            if ( !url || !( await isUrlReachable( url ) ) ) {
                console.error( 'No running API service detected.' );
                process.exit( 1 );
            }
            console.log( `Detecting an API service: [${url}].` );
            process.env.REACT_APP_PORTFORWARDER_API_BASE_URL = url;
            execSync( 'react-scripts start', { stdio: 'inherit' } );
            break;
        case 'build-gui':
            process.env.CGO_ENABLED = 0;
            execSync( 'react-scripts build && statik -f -src=build -p _statik -dest ./', { stdio: 'inherit' } );
            break;
        case 'build-app':
            process.env.CGO_ENABLED = 0;
            execSync( 'go build -C=./ -ldflags -H=windowsgui -o PortForwarder.exe', { stdio: 'inherit' } );
            break;
        case 'build-all':
            process.env.CGO_ENABLED = 0;
            execSync( 'npm run build:gui && npm run build:app', { stdio: 'inherit' } );
            break;
        default:
            console.error( `unknown action ${action}` );
            process.exit( 1 );
    }
} )();
