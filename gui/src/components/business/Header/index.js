import './index.css';
import Menu from '../Menu';
import Brand from '../Brand';

export default function Header() {
    return (
        <div className="header">
            <div className="logo">
                <Brand />
            </div>
            <div className="other">
                <Menu />
            </div>
        </div>
    );
}
