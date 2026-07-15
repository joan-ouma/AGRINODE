import React, { useState, useEffect } from 'react';
import { 
  Activity, Cpu, Wifi, Settings, Database, BarChart3,
  Thermometer, Droplets, AlertTriangle, Menu, X, Bell, Moon
} from 'lucide-react';
import { 
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, 
  ResponsiveContainer, AreaChart, Area
} from 'recharts';

function App() {
  const [data, setData] = useState([]);
  const [historyData, setHistoryData] = useState([]);
  const [currentTemp, setCurrentTemp] = useState(0);
  const [currentHumidity, setCurrentHumidity] = useState(0);
  const [isLive, setIsLive] = useState(false);
  const [error, setError] = useState(null);
  
  // Navigation State
  const [activeTab, setActiveTab] = useState('Dashboard');
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  
  // Settings State
  const [darkMode, setDarkMode] = useState(true);
  const [notifications, setNotifications] = useState(true);

  // Fetch real data from the Go API
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('/api/v1/analytics/daily');
        if (!response.ok) throw new Error('Failed to fetch data');
        
        const result = await response.json();
        
        if (result.status === 'success' && result.data) {
          const formattedData = result.data.reverse().map(item => ({
            time: new Date(item.reading_date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }),
            temperature: item.avg_temp,
            humidity: item.avg_humidity,
            moisture: item.avg_moisture,
            readings: item.total_daily_readings,
            node_name: item.node_name,
            zone: item.zone
          }));
          
          setData(formattedData);
          
          if (formattedData.length > 0) {
            const latest = formattedData[formattedData.length - 1];
            setCurrentTemp(latest.temperature.toFixed(1));
            setCurrentHumidity(latest.humidity.toFixed(1));
          }
          
          setIsLive(true);
          setError(null);
        }
      } catch (err) {
        console.error("API Error:", err);
        setError("API connection lost or unavailable.");
        setIsLive(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 10000);
    return () => clearInterval(interval);
  }, []);

  // Fetch Historical Data when Data Logs tab is active
  useEffect(() => {
    if (activeTab === 'Data Logs') {
      const fetchHistory = async () => {
        try {
          const response = await fetch('/api/history');
          if (!response.ok) throw new Error('Failed to fetch history');
          const result = await response.json();
          setHistoryData(result || []);
        } catch (err) {
          console.error("History API Error:", err);
        }
      };
      fetchHistory();
      const interval = setInterval(fetchHistory, 15000);
      return () => clearInterval(interval);
    }
  }, [activeTab]);

  const handleNavClick = (tab) => {
    setActiveTab(tab);
    setIsMobileMenuOpen(false); // Close menu on mobile after click
  };

  const renderDashboard = () => (
    <>
      <div className="grid-container">
        <div className="glass-card stat-card">
          <div className="stat-header">
            <span>Avg Temperature</span>
            <Thermometer size={20} color="var(--accent-warning)" />
          </div>
          <div className="stat-value">{currentTemp}°C</div>
          <div className="stat-footer trend-up">
            <Activity size={14} />
            <span>{isLive ? 'Live Tracking' : 'Offline'}</span>
          </div>
        </div>

        <div className="glass-card stat-card">
          <div className="stat-header">
            <span>Avg Humidity</span>
            <Droplets size={20} color="var(--accent-secondary)" />
          </div>
          <div className="stat-value">{currentHumidity}%</div>
          <div className="stat-footer trend-down">
            <Activity size={14} />
            <span>{isLive ? 'Live Tracking' : 'Offline'}</span>
          </div>
        </div>

        <div className="glass-card stat-card">
          <div className="stat-header">
            <span>Active Nodes</span>
            <Wifi size={20} color="var(--accent-primary)" />
          </div>
          <div className="stat-value">{isLive ? '1 / 1' : '0 / 1'}</div>
          <div className="stat-footer trend-up" style={{ color: isLive ? 'var(--accent-primary)' : 'var(--text-muted)' }}>
            <Activity size={14} />
            <span>{isLive ? '100% Connectivity' : 'No connection'}</span>
          </div>
        </div>

        <div className="glass-card stat-card" style={{ borderColor: 'rgba(255, 59, 48, 0.3)' }}>
          <div className="stat-header">
            <span>Anomalies</span>
            <AlertTriangle size={20} color="var(--accent-danger)" />
          </div>
          <div className="stat-value">0</div>
          <div className="stat-footer">
            <Activity size={14} color="var(--text-muted)" />
            <span style={{ color: 'var(--text-muted)' }}>No issues detected</span>
          </div>
        </div>
      </div>

      <div className="glass-card chart-container">
        <div style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between' }}>
          <h3 style={{ fontSize: '1.2rem' }}>Sensor Telemetry Stream (Daily Averages)</h3>
        </div>
        <div className="chart-wrapper" style={{ height: '320px', width: '100%' }}>
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="colorTemp" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="var(--accent-warning)" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="var(--accent-warning)" stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorHum" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="var(--accent-secondary)" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="var(--accent-secondary)" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--border-light)" vertical={false} />
              <XAxis dataKey="time" stroke="var(--text-muted)" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis stroke="var(--text-muted)" fontSize={12} tickLine={false} axisLine={false} />
              <Tooltip 
                contentStyle={{ 
                  backgroundColor: 'var(--bg-card)', 
                  border: '1px solid var(--border-light)',
                  borderRadius: '8px',
                  backdropFilter: 'blur(8px)'
                }}
                itemStyle={{ color: 'var(--text-main)' }}
              />
              <Area type="monotone" dataKey="temperature" stroke="var(--accent-warning)" fillOpacity={1} fill="url(#colorTemp)" strokeWidth={2} name="Avg Temp (°C)" />
              <Area type="monotone" dataKey="humidity" stroke="var(--accent-secondary)" fillOpacity={1} fill="url(#colorHum)" strokeWidth={2} name="Avg Humidity (%)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>
    </>
  );

  const renderDevices = () => {
    // Extract unique devices from the data stream
    const uniqueDevices = [];
    const map = new Map();
    for (const item of data) {
      if(item.node_name && !map.has(item.node_name)){
        map.set(item.node_name, true);
        uniqueDevices.push({
          name: item.node_name,
          zone: item.zone,
          status: isLive ? 'Online' : 'Offline',
          readings: item.readings
        });
      }
    }

    if(uniqueDevices.length === 0 && isLive) {
       uniqueDevices.push({ name: 'ESP8266-Alpha', zone: 'Zone 1', status: 'Online', readings: 'Live' });
    }

    return (
      <div className="grid-container">
        {uniqueDevices.map((device, i) => (
          <div key={i} className="glass-card stat-card">
            <div className="stat-header">
              <span>{device.name}</span>
              <Cpu size={20} color={device.status === 'Online' ? 'var(--accent-primary)' : 'var(--text-muted)'} />
            </div>
            <div style={{ padding: '10px 0' }}>
              <div style={{ fontSize: '0.9rem', color: 'var(--text-muted)', marginBottom: '5px' }}>Location: {device.zone}</div>
              <div style={{ fontSize: '0.9rem', color: 'var(--text-muted)' }}>Activity: {device.readings} readings today</div>
            </div>
            <div className="stat-footer" style={{ color: device.status === 'Online' ? 'var(--accent-primary)' : 'var(--accent-danger)' }}>
              <Wifi size={14} />
              <span>{device.status}</span>
            </div>
          </div>
        ))}
        {uniqueDevices.length === 0 && !isLive && (
          <div style={{ gridColumn: '1 / -1', color: 'var(--text-muted)', textAlign: 'center', padding: '40px' }}>
            No devices found. Ensure the API is connected.
          </div>
        )}
      </div>
    );
  };

  const renderDataLogs = () => (
    <div className="glass-card data-logs-card" style={{ padding: '24px' }}>
      <h3 style={{ fontSize: '1.2rem', marginBottom: '20px' }}>Historical Telemetry Logs (MongoDB)</h3>
      <div className="table-container">
        <table className="data-table">
          <thead>
            <tr>
              <th>Timestamp</th>
              <th>Temperature (°C)</th>
              <th>Humidity (%)</th>
              <th>Soil Moisture</th>
            </tr>
          </thead>
          <tbody>
            {historyData.length > 0 ? (
              historyData.map((log, i) => (
                <tr key={i}>
                  <td>{new Date(log.timestamp).toLocaleString()}</td>
                  <td>{log.temperature?.toFixed(2)}</td>
                  <td>{log.humidity?.toFixed(2)}</td>
                  <td>{log.soilMoisture}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="4" style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '40px' }}>
                  {isLive ? 'Fetching logs...' : 'No data available'}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );

  const renderSettings = () => (
    <div className="glass-card settings-card" style={{ padding: '32px', maxWidth: '600px', margin: '0 auto' }}>
      <h3 style={{ fontSize: '1.2rem', marginBottom: '30px' }}>System Preferences</h3>
      <div className="settings-list">
        <div className="setting-item">
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Moon size={20} color="var(--accent-secondary)" />
            <div>
              <div style={{ fontWeight: '500' }}>Dark Mode</div>
              <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Use dark theme across the dashboard</div>
            </div>
          </div>
          <div className={`toggle-switch ${darkMode ? 'on' : ''}`} onClick={() => setDarkMode(!darkMode)}>
            <div className="toggle-knob"></div>
          </div>
        </div>

        <div className="setting-item">
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Bell size={20} color="var(--accent-warning)" />
            <div>
              <div style={{ fontWeight: '500' }}>Push Notifications</div>
              <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Receive alerts for critical anomalies</div>
            </div>
          </div>
          <div className={`toggle-switch ${notifications ? 'on' : ''}`} onClick={() => setNotifications(!notifications)}>
            <div className="toggle-knob"></div>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="app-container">
      {/* Mobile Header */}
      <div className="mobile-header">
        <div className="brand" style={{ marginBottom: 0 }}>
          <Activity color="var(--accent-primary)" size={24} />
          <h1 className="text-gradient" style={{ fontSize: '1.2rem' }}>Agri-Node</h1>
        </div>
        <button className="menu-toggle" onClick={() => setIsMobileMenuOpen(true)}>
          <Menu size={24} />
        </button>
      </div>

      {/* Sidebar Overlay for Mobile */}
      <div 
        className={`sidebar-overlay ${isMobileMenuOpen ? 'open' : ''}`} 
        onClick={() => setIsMobileMenuOpen(false)}
      />

      {/* Sidebar Navigation */}
      <aside className={`sidebar ${isMobileMenuOpen ? 'open' : ''}`}>
        <div className="brand" style={{ display: 'flex', justifyContent: 'space-between', width: '100%' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Activity color="var(--accent-primary)" size={32} />
            <h1 className="text-gradient">Agri-Node</h1>
          </div>
          <button className="menu-toggle mobile-close-btn" onClick={() => setIsMobileMenuOpen(false)}>
             <X size={24} className="mobile-close" />
          </button>
        </div>
        
        <nav style={{ flex: 1 }}>
          <a href="#" className={`nav-link ${activeTab === 'Dashboard' ? 'active' : ''}`} onClick={(e) => { e.preventDefault(); handleNavClick('Dashboard'); }}>
            <BarChart3 />
            Dashboard
          </a>
          <a href="#" className={`nav-link ${activeTab === 'Devices' ? 'active' : ''}`} onClick={(e) => { e.preventDefault(); handleNavClick('Devices'); }}>
            <Cpu />
            Devices
          </a>
          <a href="#" className={`nav-link ${activeTab === 'Data Logs' ? 'active' : ''}`} onClick={(e) => { e.preventDefault(); handleNavClick('Data Logs'); }}>
            <Database />
            Data Logs
          </a>
          <a href="#" className={`nav-link ${activeTab === 'Settings' ? 'active' : ''}`} onClick={(e) => { e.preventDefault(); handleNavClick('Settings'); }}>
            <Settings />
            Settings
          </a>
        </nav>

        <div className="glass-card" style={{ padding: '16px', marginTop: 'auto', border: `1px solid ${isLive ? 'var(--accent-primary-dim)' : 'var(--accent-danger)'}` }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            {isLive ? (
              <div className="live-badge">System Normal</div>
            ) : (
              <div className="live-badge" style={{ backgroundColor: 'rgba(255, 59, 48, 0.1)', color: 'var(--accent-danger)' }}>Disconnected</div>
            )}
          </div>
          <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginTop: '8px' }}>
            {isLive ? 'API service connected to PostgreSQL.' : 'Failed to reach API gateway.'}
          </p>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="main-content">
        <header className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '40px', flexWrap: 'wrap', gap: '16px' }}>
          <div>
            <h2 style={{ fontSize: '1.8rem', marginBottom: '8px' }}>{activeTab}</h2>
            <p style={{ color: 'var(--text-muted)' }}>
              {activeTab === 'Dashboard' && 'Real-time sensor data & system analytics'}
              {activeTab === 'Devices' && 'Manage and monitor connected edge nodes'}
              {activeTab === 'Data Logs' && 'Historical event logs and raw telemetry records'}
              {activeTab === 'Settings' && 'Configure dashboard preferences and alerts'}
            </p>
          </div>
          {isLive ? (
            <div className="live-badge">Live API Connected</div>
          ) : (
            <div className="live-badge" style={{ backgroundColor: 'rgba(255, 59, 48, 0.1)', color: 'var(--accent-danger)' }}>
              API Offline
            </div>
          )}
        </header>

        {error && (
          <div className="glass-card" style={{ padding: '16px', marginBottom: '24px', borderLeft: '4px solid var(--accent-danger)' }}>
            <p style={{ color: 'var(--accent-danger)' }}>{error}</p>
          </div>
        )}

        {/* Dynamic Content Rendering */}
        {activeTab === 'Dashboard' && renderDashboard()}
        {activeTab === 'Devices' && renderDevices()}
        {activeTab === 'Data Logs' && renderDataLogs()}
        {activeTab === 'Settings' && renderSettings()}

      </main>
    </div>
  );
}

export default App;
