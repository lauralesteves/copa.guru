import { useMemo, useRef } from 'react';
import { Canvas, useFrame, useThree } from '@react-three/fiber';
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

// Custom shader for a stylized globe with grid lines and gradient
const globeVertexShader = `
  varying vec3 vNormal;
  varying vec3 vPosition;
  varying vec2 vUv;

  void main() {
    vNormal = normalize(normalMatrix * normal);
    vPosition = position;
    vUv = uv;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
  }
`;

const globeFragmentShader = `
  varying vec3 vNormal;
  varying vec3 vPosition;
  varying vec2 vUv;

  uniform vec3 uLightDir;

  void main() {
    // Base color: dark purple ocean
    vec3 ocean = vec3(0.04, 0.005, 0.08);
    // Land color: dark teal/green
    vec3 land = vec3(0.06, 0.12, 0.1);

    // Simple procedural continents using spherical coords
    float lat = asin(normalize(vPosition).y);
    float lng = atan(normalize(vPosition).z, normalize(vPosition).x);

    // Create "continents" using layered sine waves
    float continent = 0.0;
    continent += sin(lat * 3.0 + 0.5) * sin(lng * 2.0 - 0.3) * 0.5;
    continent += sin(lat * 5.0 - 1.0) * sin(lng * 3.0 + 1.0) * 0.3;
    continent += sin(lat * 7.0 + 2.0) * sin(lng * 5.0 - 0.5) * 0.2;
    continent = smoothstep(0.1, 0.3, continent);

    vec3 baseColor = mix(ocean, land, continent);

    // Grid lines (latitude and longitude)
    float latLine = abs(fract(lat * 5.73) - 0.5);
    float lngLine = abs(fract(lng * 5.73) - 0.5);
    float grid = 1.0 - smoothstep(0.01, 0.03, min(latLine, lngLine));
    vec3 gridColor = vec3(0.83, 0.66, 0.26); // gold
    baseColor = mix(baseColor, gridColor, grid * 0.15);

    // Equator highlight
    float equator = 1.0 - smoothstep(0.01, 0.04, abs(lat));
    baseColor = mix(baseColor, gridColor, equator * 0.2);

    // Lighting
    float diffuse = max(dot(vNormal, normalize(uLightDir)), 0.0);
    float ambient = 0.3;
    float light = ambient + diffuse * 0.7;

    // Fresnel rim
    vec3 viewDir = normalize(cameraPosition - vPosition);
    float fresnel = pow(1.0 - max(dot(vNormal, viewDir), 0.0), 3.0);
    vec3 rimColor = vec3(0.1, 0.02, 0.2);

    vec3 finalColor = baseColor * light + rimColor * fresnel * 0.6;

    gl_FragColor = vec4(finalColor, 0.95);
  }
`;

function GlobeMesh() {
  const shaderRef = useRef<THREE.ShaderMaterial>(null);

  const uniforms = useMemo(
    () => ({
      uLightDir: { value: new THREE.Vector3(5, 3, 5).normalize() },
    }),
    [],
  );

  return (
    <mesh>
      <sphereGeometry args={[GLOBE_RADIUS, 64, 64]} />
      <shaderMaterial
        ref={shaderRef}
        vertexShader={globeVertexShader}
        fragmentShader={globeFragmentShader}
        uniforms={uniforms}
        transparent
      />
    </mesh>
  );
}

function Globe({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  const globeRef = useRef<THREE.Group>(null);
  const { raycaster, pointer, camera } = useThree();
  const pointMeshesRef = useRef<THREE.Mesh[]>([]);

  useFrame((_, delta) => {
    if (globeRef.current) {
      globeRef.current.rotation.y += delta * 0.08;
    }
  });

  useFrame(() => {
    if (pointMeshesRef.current.length === 0) return;
    raycaster.setFromCamera(pointer, camera);
    const intersects = raycaster.intersectObjects(pointMeshesRef.current);
    if (intersects.length > 0) {
      const idx = pointMeshesRef.current.indexOf(intersects[0].object as THREE.Mesh);
      if (idx >= 0) {
        onPointHover(pointsData[idx]);
        return;
      }
    }
    onPointHover(null);
  });

  return (
    <group ref={globeRef}>
      <GlobeMesh />

      {/* Atmosphere glow */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.12, 32, 32]} />
        <meshBasicMaterial
          color="#2d0a54"
          transparent
          opacity={0.12}
          side={THREE.BackSide}
        />
      </mesh>

      {/* Second atmosphere layer */}
      <mesh>
        <sphereGeometry args={[GLOBE_RADIUS + 0.25, 32, 32]} />
        <meshBasicMaterial
          color="#d4a843"
          transparent
          opacity={0.03}
          side={THREE.BackSide}
        />
      </mesh>

      {/* Points */}
      {pointsData.map((point, i) => {
        const pos = latLngToVector3(point.lat, point.lng, GLOBE_RADIUS);
        const dir = pos.clone().normalize();
        const height = point.size * 0.5;

        return (
          <mesh
            key={point.code}
            position={[
              pos.x + dir.x * height / 2,
              pos.y + dir.y * height / 2,
              pos.z + dir.z * height / 2,
            ]}
            ref={(el) => {
              if (el) pointMeshesRef.current[i] = el;
            }}
            lookAt={[0, 0, 0]}
          >
            <cylinderGeometry args={[0.025, 0.04, height, 8]} />
            <meshStandardMaterial
              color={point.color}
              emissive={point.color}
              emissiveIntensity={0.6}
            />
          </mesh>
        );
      })}

      {/* Arcs */}
      {arcsData.map((arc, i) => {
        const start = latLngToVector3(arc.startLat, arc.startLng, GLOBE_RADIUS + 0.01);
        const end = latLngToVector3(arc.endLat, arc.endLng, GLOBE_RADIUS + 0.01);
        const mid = start.clone().add(end).multiplyScalar(0.5).normalize().multiplyScalar(GLOBE_RADIUS * 1.35);

        const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
        const points = curve.getPoints(32);
        const geometry = new THREE.BufferGeometry().setFromPoints(points);

        return (
          <primitive
            key={`arc-${i}-${arc.startLat.toFixed(1)}-${arc.endLat.toFixed(1)}`}
            object={new THREE.Line(geometry, new THREE.LineBasicMaterial({ color: arc.color, transparent: true, opacity: 0.18 }))}
          />
        );
      })}
    </group>
  );
}

export default function GlobeCanvas({ pointsData, arcsData, onPointHover }: GlobeCanvasProps) {
  return (
    <Canvas
      camera={{ position: [0, 0, 5.5], fov: 40 }}
      dpr={[1, 1.5]}
      gl={{ alpha: true, antialias: true }}
      style={{ background: 'transparent' }}
    >
      <ambientLight intensity={0.5} />
      <directionalLight position={[5, 3, 5]} intensity={0.8} />
      <pointLight position={[-5, -3, -5]} intensity={0.3} color="#d4a843" />
      <Globe pointsData={pointsData} arcsData={arcsData} onPointHover={onPointHover} />
      <OrbitControls
        enableZoom={false}
        enablePan={false}
        rotateSpeed={0.5}
      />
    </Canvas>
  );
}
