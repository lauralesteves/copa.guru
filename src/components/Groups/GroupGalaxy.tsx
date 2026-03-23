import { useRef, useMemo } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Text, OrbitControls } from '@react-three/drei';
import * as THREE from 'three';
import type { GroupStanding } from '../../types/worldcup';
import { getTeamByCode } from '../../data/teams';

interface GroupGalaxyProps {
  teamCodes: string[];
  standings: GroupStanding[];
  groupName: string;
}

function OrbitingTeam({
  teamCode,
  orbitRadius,
  speed,
  size,
  delay,
}: {
  teamCode: string;
  orbitRadius: number;
  speed: number;
  size: number;
  delay: number;
}) {
  const groupRef = useRef<THREE.Group>(null);
  const team = getTeamByCode(teamCode);
  const color = team?.primaryColor ?? '#d4a843';

  useFrame((state) => {
    if (groupRef.current) {
      const t = state.clock.elapsedTime * speed + delay;
      groupRef.current.position.x = Math.cos(t) * orbitRadius;
      groupRef.current.position.z = Math.sin(t) * orbitRadius;
      groupRef.current.position.y = Math.sin(t * 0.5) * 0.15;
    }
  });

  return (
    <>
      {/* Orbit ring */}
      <mesh rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[orbitRadius, 0.005, 8, 64]} />
        <meshBasicMaterial color="#ffffff" transparent opacity={0.05} />
      </mesh>

      <group ref={groupRef}>
        <mesh>
          <sphereGeometry args={[size, 16, 16]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={0.4}
            roughness={0.4}
            metalness={0.6}
          />
        </mesh>
        <Text
          position={[0, -size - 0.08, 0]}
          fontSize={0.06}
          color="white"
          anchorX="center"
          anchorY="top"
          font="https://fonts.gstatic.com/s/inter/v18/UcCO3FwrK3iLTeHuS_nVMrMxCp50SjIw2boKoduKmMEVuLyfAZ9hiA.woff2"
        >
          {team?.name ?? teamCode}
        </Text>
      </group>
    </>
  );
}

function GalaxyScene({ teamCodes, standings, groupName }: GroupGalaxyProps) {
  const validTeams = teamCodes.filter((c) => !c.startsWith('TBD'));

  const teamData = useMemo(() => {
    return validTeams.map((code, i) => {
      const standing = standings.find((s) => s.teamCode === code);
      const pos = standings.findIndex((s) => s.teamCode === code);
      const points = standing?.points ?? 0;
      return {
        code,
        orbitRadius: 0.4 + (pos >= 0 ? pos : i) * 0.35,
        speed: 0.6 - (pos >= 0 ? pos : i) * 0.1,
        size: 0.08 + Math.min(points, 9) * 0.015,
        delay: i * (Math.PI / 2),
      };
    });
  }, [validTeams, standings]);

  return (
    <>
      <ambientLight intensity={0.4} />
      <pointLight position={[0, 0, 0]} intensity={1.5} color="#d4a843" distance={5} />

      {/* Central star */}
      <mesh>
        <sphereGeometry args={[0.12, 16, 16]} />
        <meshStandardMaterial
          color="#d4a843"
          emissive="#d4a843"
          emissiveIntensity={1}
        />
      </mesh>

      {/* Group label */}
      <Text
        position={[0, 0.25, 0]}
        fontSize={0.12}
        color="#d4a843"
        anchorX="center"
        font="https://fonts.gstatic.com/s/inter/v18/UcCO3FwrK3iLTeHuS_nVMrMxCp50SjIw2boKoduKmMEVuLyfAZ9hiA.woff2"
        fontWeight={700}
      >
        {groupName}
      </Text>

      {teamData.map((t) => (
        <OrbitingTeam key={t.code} teamCode={t.code} {...t} />
      ))}

      <OrbitControls enableZoom={false} enablePan={false} autoRotate autoRotateSpeed={0.5} />
    </>
  );
}

export function GroupGalaxy({ teamCodes, standings, groupName }: GroupGalaxyProps) {
  return (
    <div className="w-full h-[200px]">
      <Canvas
        camera={{ position: [0, 1.2, 2.2], fov: 45 }}
        dpr={[1, 1.5]}
        gl={{ alpha: true, antialias: true }}
        style={{ background: 'transparent' }}
      >
        <GalaxyScene teamCodes={teamCodes} standings={standings} groupName={groupName} />
      </Canvas>
    </div>
  );
}
