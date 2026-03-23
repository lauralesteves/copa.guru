import { useRef, useMemo } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Float, Text, OrbitControls } from '@react-three/drei';
import * as THREE from 'three';
import { getTeamByCode } from '../../data/teams';
import type { GroupName, GroupStanding } from '../../types/worldcup';

interface Bracket3DProps {
  allGroupStandings: Map<GroupName, GroupStanding[]>;
}

interface Node3D {
  id: string;
  teamCode: string | null;
  label: string;
  position: [number, number, number];
  size: number;
  stage: number; // 0=r32, 1=r16, 2=qf, 3=sf, 4=final
}

interface Link3D {
  from: [number, number, number];
  to: [number, number, number];
}

function buildNodes(
  allGroupStandings: Map<GroupName, GroupStanding[]>,
): { nodes: Node3D[]; links: Link3D[] } {
  const nodes: Node3D[] = [];
  const links: Link3D[] = [];

  const groupPairs: [GroupName, GroupName][] = [
    ['A', 'B'], ['C', 'D'], ['E', 'F'], ['G', 'H'], ['I', 'J'], ['K', 'L'],
  ];

  function getTeam(group: GroupName, pos: number): string | null {
    const s = allGroupStandings.get(group);
    if (!s?.[pos] || s[pos].points === 0) return null;
    return s[pos].teamCode;
  }

  // R32 (24 teams: 1st and 2nd from each group)
  const r32Positions: [number, number, number][] = [];
  let r32i = 0;
  for (const [g1, g2] of groupPairs) {
    const teams = [
      { code: getTeam(g1, 0), label: `1${g1}` },
      { code: getTeam(g2, 1), label: `2${g2}` },
      { code: getTeam(g2, 0), label: `1${g2}` },
      { code: getTeam(g1, 1), label: `2${g1}` },
    ];
    for (const t of teams) {
      const y = (r32i - 11.5) * 0.55;
      const pos: [number, number, number] = [-6, y, 0];
      r32Positions.push(pos);
      nodes.push({ id: `r32-${r32i}`, teamCode: t.code, label: t.label, position: pos, size: 0.15, stage: 0 });
      r32i++;
    }
  }

  // R16 (12 nodes)
  const r16Positions: [number, number, number][] = [];
  for (let i = 0; i < 12; i++) {
    const y = (i - 5.5) * 1.1;
    const pos: [number, number, number] = [-3, y, -1.5];
    r16Positions.push(pos);
    nodes.push({ id: `r16-${i}`, teamCode: null, label: 'Oitavas', position: pos, size: 0.2, stage: 1 });
    if (i < r32Positions.length / 2) {
      links.push({ from: r32Positions[i * 2], to: pos });
      links.push({ from: r32Positions[i * 2 + 1], to: pos });
    }
  }

  // QF (6 nodes)
  const qfPositions: [number, number, number][] = [];
  for (let i = 0; i < 6; i++) {
    const y = (i - 2.5) * 2;
    const pos: [number, number, number] = [0, y, -3];
    qfPositions.push(pos);
    nodes.push({ id: `qf-${i}`, teamCode: null, label: 'Quartas', position: pos, size: 0.25, stage: 2 });
    if (i < r16Positions.length / 2) {
      links.push({ from: r16Positions[i * 2], to: pos });
      links.push({ from: r16Positions[i * 2 + 1], to: pos });
    }
  }

  // SF (3 nodes — simplified)
  const sfPositions: [number, number, number][] = [];
  for (let i = 0; i < 3; i++) {
    const y = (i - 1) * 3;
    const pos: [number, number, number] = [3, y, -4.5];
    sfPositions.push(pos);
    nodes.push({ id: `sf-${i}`, teamCode: null, label: 'Semi', position: pos, size: 0.3, stage: 3 });
    if (i < qfPositions.length / 2) {
      links.push({ from: qfPositions[i * 2], to: pos });
      links.push({ from: qfPositions[i * 2 + 1], to: pos });
    }
  }

  // Final
  const finalPos: [number, number, number] = [6, 0, -6];
  nodes.push({ id: 'final', teamCode: null, label: 'FINAL', position: finalPos, size: 0.45, stage: 4 });
  if (sfPositions.length >= 2) {
    links.push({ from: sfPositions[0], to: finalPos });
    links.push({ from: sfPositions[sfPositions.length - 1], to: finalPos });
  }

  return { nodes, links };
}

function TubeConnection({ from, to }: { from: [number, number, number]; to: [number, number, number] }) {
  const geometry = useMemo(() => {
    const start = new THREE.Vector3(...from);
    const end = new THREE.Vector3(...to);
    const mid = start.clone().add(end).multiplyScalar(0.5);
    mid.z -= 0.5;
    const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
    return new THREE.TubeGeometry(curve, 16, 0.015, 6, false);
  }, [from, to]);

  return (
    <mesh geometry={geometry}>
      <meshBasicMaterial color="#d4a843" transparent opacity={0.15} />
    </mesh>
  );
}

function BracketNode({ node }: { node: Node3D }) {
  const meshRef = useRef<THREE.Mesh>(null);
  const team = node.teamCode ? getTeamByCode(node.teamCode) : null;
  const color = team?.primaryColor ?? (node.stage === 4 ? '#d4a843' : '#1a0533');

  useFrame((state) => {
    if (meshRef.current && node.stage === 4) {
      meshRef.current.scale.setScalar(1 + Math.sin(state.clock.elapsedTime * 2) * 0.05);
    }
  });

  return (
    <Float speed={1} rotationIntensity={0} floatIntensity={node.stage === 4 ? 0.3 : 0.1}>
      <group position={node.position}>
        <mesh ref={meshRef}>
          <sphereGeometry args={[node.size, 16, 16]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={node.stage === 4 ? 0.6 : 0.2}
            transparent
            opacity={node.teamCode ? 0.9 : 0.4}
            roughness={0.3}
            metalness={0.5}
          />
        </mesh>

        {/* Team label */}
        <Text
          position={[0, -node.size - 0.12, 0]}
          fontSize={0.1}
          color={node.teamCode ? '#ffffff' : 'rgba(255,255,255,0.3)'}
          anchorX="center"
          anchorY="top"
          font="https://fonts.gstatic.com/s/inter/v18/UcCO3FwrK3iLTeHuS_nVMrMxCp50SjIw2boKoduKmMEVuLyfAZ9hiA.woff2"
        >
          {node.teamCode ? (team?.name ?? node.teamCode) : node.label}
        </Text>
      </group>
    </Float>
  );
}

function Scene({ allGroupStandings }: Bracket3DProps) {
  const { nodes, links } = useMemo(() => buildNodes(allGroupStandings), [allGroupStandings]);

  return (
    <>
      <ambientLight intensity={0.3} />
      <directionalLight position={[10, 5, 5]} intensity={0.8} />
      <pointLight position={[6, 0, -6]} intensity={1} color="#d4a843" distance={10} />

      {links.map((link, i) => (
        <TubeConnection key={`link-${i}-${link.from[0]}-${link.to[0]}`} from={link.from} to={link.to} />
      ))}

      {nodes.map((node) => (
        <BracketNode key={node.id} node={node} />
      ))}

      <OrbitControls
        enableZoom
        enablePan
        minDistance={3}
        maxDistance={20}
        autoRotate
        autoRotateSpeed={0.3}
      />
    </>
  );
}

export function Bracket3D({ allGroupStandings }: Bracket3DProps) {
  return (
    <div className="w-full h-[500px] rounded-xl overflow-hidden border border-white/5">
      <Canvas
        camera={{ position: [0, 0, 12], fov: 50 }}
        dpr={[1, 1.5]}
        gl={{ alpha: true, antialias: true }}
        style={{ background: 'transparent' }}
      >
        <Scene allGroupStandings={allGroupStandings} />
      </Canvas>
    </div>
  );
}
