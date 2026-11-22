import { forwardRef, useEffect, useMemo, useRef, useState } from 'react';
import { Canvas, useFrame, useLoader, useThree } from '@react-three/fiber';
import { OrbitControls } from '@react-three/drei';
import * as THREE from 'three';

interface PointData {
  lat: number;
  lng: number;
  code: string;
  name: string;
  color: string;
  group: string;
  size: number;
  confederation: string;
}

interface ArcData {
  startLat: number;
  startLng: number;
  endLat: number;
  endLng: number;
  color: string;
}

interface GlobeCanvasProps {
  pointsData: PointData[];
  arcsData: ArcData[];
  onPointHover: (point: PointData | null) => void;
}

const GLOBE_RADIUS = 2;

function latLngToVector3(lat: number, lng: number, radius: number): THREE.Vector3 {
  const phi = (90 - lat) * (Math.PI / 180);
  const theta = (lng + 180) * (Math.PI / 180);
  return new THREE.Vector3(
    -radius * Math.sin(phi) * Math.cos(theta),
    radius * Math.cos(phi),
    radius * Math.sin(phi) * Math.sin(theta),
  );
}

function EarthSphere() {
  const texture = useLoader(
    THREE.TextureLoader,
    'https://unpkg.com/three-globe/example/img/earth-blue-marble.jpg',
  );
  const bumpMap = useLoader(
    THREE.TextureLoader,
    'https://unpkg.com/three-globe/example/img/earth-topology.png',
  );

  return (
    <mesh>
      <sphereGeometry args={[GLOBE_RADIUS, 64, 64]} />
      <meshStandardMaterial
        map={texture}
        bumpMap={bumpMap}
        bumpScale={0.03}
        roughness={0.85}
        metalness={0.1}
      />
    </mesh>
  );
}

const FLAG_ALTITUDE = 0.02;
const _worldPos = new THREE.Vector3();
const _toCamera = new THREE.Vector3();

function flagPosition(point: PointData): THREE.Vector3 {
  return latLngToVector3(point.lat, point.lng, GLOBE_RADIUS + FLAG_ALTITUDE);
}

const FlagSprite = forwardRef<THREE.Sprite, { point: PointData }>(function FlagSprite({ point }, ref) {
  const texture = useLoader(THREE.TextureLoader, `/images/flags/96/${point.code}.png`);
  const position = useMemo(() => flagPosition(point), [point]);
  const spriteScale = point.size > 0.5 ? 0.22 : 0.16;

  return (
    <sprite ref={ref} position={position} scale={[spriteScale, spriteScale * 0.7, 1]} renderOrder={1}>
      <spriteMaterial map={texture} transparent opacity={0.95} depthTest={false} depthWrite={false} />
    </sprite>
  );
});

function Globe({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  const globeRef = useRef<THREE.Group>(null);
  const { raycaster, pointer, camera } = useThree();
  const spriteRefs = useRef<(THREE.Sprite | null)[]>([]);

  useFrame((_, delta) => {
    if (globeRef.current) {
      globeRef.current.rotation.y += delta * 0.06;
    }

    // Flags ignore the depth buffer so they sit cleanly on top of their country;
    // hide the ones on the far side of the globe instead.
    for (const sprite of spriteRefs.current) {
      if (!sprite) continue;
      sprite.getWorldPosition(_worldPos);
      _toCamera.copy(camera.position).sub(_worldPos);
      sprite.visible = _worldPos.normalize().dot(_toCamera.normalize()) > 0.05;
    }

    // Raycasting for hover on visible flags
    const targets = spriteRefs.current.filter((s): s is THREE.Sprite => !!s && s.visible);
    raycaster.setFromCamera(pointer, camera);
    const hit = raycaster.intersectObjects(targets)[0];
    const idx = hit ? spriteRefs.current.indexOf(hit.object as THREE.Sprite) : -1;
    onPointHover(idx >= 0 ? pointsData[idx] : null);
  });

  return (
    <group ref={globeRef}>
      <EarthSphere />

      {/* Atmosphere glow */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.08, 32, 32]} />
        <meshBasicMaterial
          color="#4466aa"
          transparent
          opacity={0.06}
          side={THREE.BackSide}
        />
      </mesh>
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.2, 32, 32]} />
        <meshBasicMaterial
          color="#121826"
          transparent
          opacity={0.1}
          side={THREE.BackSide}
        />
      </mesh>

      {/* Flag sprites */}
      {pointsData.map((point, i) => (
        <FlagSprite
          key={point.code}
          point={point}
          ref={(el) => {
            spriteRefs.current[i] = el;
          }}
        />
      ))}

      {/* Arcs */}
      {arcsData.map((arc, i) => {
        const start = latLngToVector3(arc.startLat, arc.startLng, GLOBE_RADIUS + 0.01);
        const end = latLngToVector3(arc.endLat, arc.endLng, GLOBE_RADIUS + 0.01);
        const mid = start.clone().add(end).multiplyScalar(0.5).normalize().multiplyScalar(GLOBE_RADIUS * 1.3);
        const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
        const points = curve.getPoints(32);
        const geometry = new THREE.BufferGeometry().setFromPoints(points);

        return (
          <primitive
            key={`arc-${i}-${arc.startLat.toFixed(1)}-${arc.endLat.toFixed(1)}`}
            object={new THREE.Line(geometry, new THREE.LineBasicMaterial({ color: arc.color, transparent: true, opacity: 0.15 }))}
          />
        );
      })}
    </group>
  );
}

export default function GlobeCanvas({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => setVisible(entry.isIntersecting),
      { threshold: 0.1 },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  return (
    <div ref={containerRef} className="w-full h-full">
      {visible && (
        <Canvas
          camera={{ position: [0, 0, 7], fov: 40 }}
          dpr={[1, 1.5]}
          gl={{ alpha: true, antialias: true, powerPreference: 'low-power' }}
          style={{ background: 'transparent' }}
        >
          <ambientLight intensity={0.8} />
          <directionalLight position={[5, 3, 5]} intensity={1} />
          <pointLight position={[-5, -3, -5]} intensity={0.3} color="#00ff9c" />
          <Globe pointsData={pointsData} arcsData={arcsData} onPointHover={onPointHover} />
          <OrbitControls
            enableZoom={false}
            enablePan={false}
            rotateSpeed={0.5}
          />
        </Canvas>
      )}
    </div>
  );
}
