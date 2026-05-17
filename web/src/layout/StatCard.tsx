type StatCardProps = {
  label: string;
  value: string | number;
  icon?: string;
  tone?: "blue" | "green" | "orange" | "red" | "purple";
};

export function StatCard({ label, value, icon = "▣", tone = "blue" }: StatCardProps) {
  return (
    <div className={`stat-card ${tone}`}>
      <div className="stat-icon">{icon}</div>
      <strong>{value}</strong>
      <span>{label}</span>
    </div>
  );
}