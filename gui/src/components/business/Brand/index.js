import './index.css';
import icon from './icon.png';

export default function Brand() {
    return (
        <div className="brand"><img src={icon} alt="PortForwarder" /><span>PortForwarder</span></div>
    );
}
