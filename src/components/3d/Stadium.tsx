import { useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls } from '@react-three/drei';
import type * as THREE from 'three';

interface StadiumProps {
  team1Color?: string;
  team2Color?: string;
}

function StadiumScene({ team1Color = '#006847', team2Color = '#002868' }: StadiumProps) {
  const stadiumRef = useRef<THREE.Group>(null);

  useFrame((state) => {
    if (stadiumRef.current) {
      stadiumRef.current.rotation.y = Math.sin(state.clock.elapsedTime * 0.15) * 0.15;
    }
  });

  return (
    <>
      <ambientLight intensity={0.4} />
      <directionalLight position={[5, 8, 5]} intensity={0.7} castShadow />
      <pointLight position={[-3, 6, -3]} intensity={0.3} color="#f0d078" />

      <group ref={stadiumRef}>
        {/* Field */}
        <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0, 0]}>
          <planeGeometry args={[4, 2.6]} />
          <meshStandardMaterial color="#1a6b30" roughness={0.9} />
        </mesh>

        {/* Field lines - center */}
        <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0.001, 0]}>
          <ringGeometry args={[0.25, 0.27, 32]} />
          <meshBasicMaterial color="white" transparent opacity={0.4} />
        </mesh>

        {/* Center line */}
        <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0.001, 0]}>
          <planeGeometry args={[0.02, 2.6]} />
          <meshBasicMaterial color="white" transparent opacity={0.4} />
        </mesh>

        {/* Stands - left side */}
        <mesh position={[-2.2, 0.5, 0]} rotation={[0, 0, 0.3]}>
          <boxGeometry args={[0.5, 1.2, 2.8]} />
          <meshStandardMaterial color={team1Color} roughness={0.7} />
        </mesh>

        {/* Stands - right side */}
        <mesh position={[2.2, 0.5, 0]} rotation={[0, 0, -0.3]}>
          <boxGeometry args={[0.5, 1.2, 2.8]} />
          <meshStandardMaterial color={team2Color} roughness={0.7} />
        </mesh>

        {/* Stands - back */}
        <mesh position={[0, 0.4, -1.5]} rotation={[-0.2, 0, 0]}>
          <boxGeometry args={[4, 0.8, 0.4]} />
          <meshStandardMaterial color="#1a0533" roughness={0.7} />
        </mesh>

        {/* Stands - front */}
        <mesh position={[0, 0.4, 1.5]} rotation={[0.2, 0, 0]}>
          <boxGeometry args={[4, 0.8, 0.4]} />
          <meshStandardMaterial color="#1a0533" roughness={0.7} />
        </mesh>

        {/* Floodlight towers */}
        {[
          [-1.8, 0, -1.2],
          [1.8, 0, -1.2],
          [-1.8, 0, 1.2],
          [1.8, 0, 1.2],
        ].map((pos) => (
          <group key={`tower-${pos[0]}-${pos[2]}`} position={pos as [number, number, number]}>
            <mesh position={[0, 1.2, 0]}>
              <cylinderGeometry args={[0.03, 0.04, 2.4, 6]} />
              <meshStandardMaterial color="#444" metalness={0.8} roughness={0.3} />
            </mesh>
            <pointLight
              position={[0, 2.5, 0]}
              intensity={0.6}
              color="#fff5e0"
              distance={6}
            />
          </group>
        ))}

        {/* Goals */}
        {[-1.95, 1.95].map((x) => (
          <group key={`goal-${x}`} position={[x, 0.2, 0]}>
            <mesh>
              <boxGeometry args={[0.05, 0.4, 0.6]} />
              <meshBasicMaterial color="white" transparent opacity={0.5} wireframe />
            </mesh>
          </group>
        ))}
      </group>

      <OrbitControls
        enableZoom={false}
        enablePan={false}
        autoRotate
        autoRotateSpeed={0.4}
        maxPolarAngle={Math.PI / 2.5}
        minPolarAngle={Math.PI / 6}
      />
    </>
  );
}

export function Stadium({ team1Color, team2Color }: StadiumProps) {
  return (
    <div className="w-full h-full min-h-[200px]">
      <Canvas
        camera={{ position: [3, 3, 4], fov: 40 }}
        dpr={[1, 1.5]}
        gl={{ alpha: true, antialias: true }}
        style={{ background: 'transparent' }}
      >
        <StadiumScene team1Color={team1Color} team2Color={team2Color} />
      </Canvas>
    </div>
  );
}
